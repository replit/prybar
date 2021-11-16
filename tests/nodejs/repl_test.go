package nodejs

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/creack/pty"
	"golang.org/x/sys/unix"
)

// There's a good chance parts of this break if we upgrade node.
// This is probably just a slight change to escape codes which human eyes won't see.
// If it does happen, all you need to do to fix it is:
//
// 1. open a repl the same way that the test would.
// 2. simulate the  test.  Make sure everything looks like it's suppposed to.
// 3. If everything looks OK, it's probably fine to just copy + paste the output from "Received: "
//    directly into the test's expected output.
//
// Some of the particular things to check are:
// - prompt / tab suggestions (though this isn't currently tested, it hopefully will be in the future)
// - preview (should evaluate any expression which can't possibly have any side effects and display a truncated version of it below the prompt)
// - results of expressions - every line should be evaulated, with the result of that line being printed to stdout
// - any errors should be displayed in red.

var prybarAssetsPath string

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	prybarAssetsPath = path.Join(cwd, "../../prybar_assets/")

	fmt.Println(prybarAssetsPath)
}

func getBasicArgs(isQuiet, isInteractive bool) []string {
	if isQuiet && isInteractive {
		return []string{"-q", "-i"}
	} else if isQuiet {
		return []string{"-q"}
	} else if isInteractive {
		return []string{"-i"}
	}

	return []string{}
}

func getCodeArgs(
	isQuiet,
	isInteractive bool,
	code string,
) []string {
	return append(getBasicArgs(isQuiet, isInteractive), "-c", code)
}

func getFileArgs(t *testing.T, isQuiet, isInteractive bool, fileContent string) []string {
	file, err := os.CreateTemp("", "prybar_test_*.js")

	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		os.Remove(file.Name())
	})

	file.WriteString(fileContent)

	return append(getBasicArgs(isQuiet, isInteractive), file.Name())

}

type testCase struct {
	name     string
	input    string
	expected string
}

type testInfo struct {
	name   string
	prefix string
	args   []string
}

func trimNewlines(str string) string {

	return str[1 : len(str)-1]
}

// used to format a string in the format of:
// ```go
//`
//<input>
//`
// into a string with '\\x1b' translated to '\x1b' and without the leading + trailing newlines.
// only needed for multiline strings.
func unescapeExpectedOutput(str string) string {
	str = trimNewlines(str)
	str = strings.ReplaceAll(str, "\\r", "\r")
	return strings.ReplaceAll(str, "\\x1b", "\x1b")
}

func escape(str string) string {
	return strings.ReplaceAll(str, "\x1b", "\\x1b")
}

type delayedReader struct {
	reader io.Reader
	once   sync.Once
}

func (r *delayedReader) Read(buf []byte) (int, error) {
	r.once.Do(func() {
		time.Sleep(50 * time.Millisecond)
	})

	return r.reader.Read(buf)
}

func (tc testCase) run(t *testing.T, info testInfo) {
	t.Parallel()
	input := bytes.NewBufferString(tc.input)
	output := &bytes.Buffer{}

	cmd := exec.Command("../../prybar-nodejs", info.args...)

	cmd.Env = append(os.Environ(), "PRYBAR_ASSETS_DIR="+prybarAssetsPath)

	ppty, tty, err := pty.Open()

	if err != nil {
		t.Fatal(err)
	}

	defer tty.Close()
	defer ppty.Close()

	// ensure that we get a constant size for the pty.
	pty.Setsize(ppty, &pty.Winsize{
		Cols: 100,
		Rows: 100,
	})

	cmd.Stdin = tty
	cmd.Stdout = tty
	cmd.Stderr = tty

	term, err := unix.IoctlGetTermios(int(tty.Fd()), unix.TCGETS)
	if err != nil {
		t.Fatal(err)
	}

	// disable local echo, the repl does this manually and we don't want to echo input twice.
	// Without this, the input shows up twice, once when we send it (before the proc starts)
	// and again later when the proc reads it from buffered input.
	term.Lflag &^= unix.ECHO

	if err := unix.IoctlSetTermios(int(tty.Fd()), unix.TCSETS, term); err != nil {
		t.Fatal(err)
	}

	if _, err := io.Copy(ppty, input); err != nil {
		panic(err)
	}
	go io.Copy(output, ppty)

	var killLock sync.Mutex
	didKill := false

	time.AfterFunc(time.Second, func() {
		killLock.Lock()
		defer killLock.Unlock()
		didKill = true
		cmd.Process.Kill()
	})

	if err := cmd.Run(); err != nil {

		// the cmd failed to start
		if cmd.ProcessState == nil {
			t.Fatal(err)
		}

		exitCode := cmd.ProcessState.ExitCode()

		killLock.Lock()
		killed := didKill
		killLock.Unlock()

		// we didn't exit with 0 (ok) or 137 (killed with SIGKILL)
		if exitCode != 0 && !killed {
			t.Fatal(err)
		}
	}

	// ignore carraige returns
	outStr := strings.ReplaceAll(output.String(), "\r", "")
	outStr = strings.TrimPrefix(outStr, info.prefix)

	if outStr != tc.expected {
		t.Log("Received:", escape(outStr))
		t.Log("Expected:", escape(tc.expected))

		t.Fail()
	}
}

func testPreview(t *testing.T, info testInfo) {
	t.Parallel()

	for _, tc := range [...]testCase{
		{
			name:  "number",
			input: "3",
			// (prompt) `3`
			expected: "\x1b[1G\x1b[0J--> \x1b[5G3",
		},
		{
			name:  "number + number",
			input: "5 + 10",
			// (prompt)  `5 + 10`
			// (preview) `15`
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5G5 + 10
\x1b[90m15\x1b[39m\x1b[11G\x1b[1A
`),
		},
		{
			name:  "module variable + number",
			input: "a + 3",
			// (prompt)  `a + 3`
			// (preview) `8`
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5Ga + 3
\x1b[90m8\x1b[39m\x1b[10G\x1b[1A
`),
		},
		{
			name:  "non-existant variable",
			input: "notAVariable",
			// (prompt)  `notAVariable`
			expected: "\x1b[1G\x1b[0J--> \x1b[5GnotAVariable",
		},
		{
			name: "repl-scoped variable",
			input: trimNewlines(`
const replScopedVariable = 'Bob'
replScopedVariable
`),
			// (prompt)  `const replScopedVariable = 'Bob'`
			// (result)  `undefined`
			// (prompt)  `replScopedVariable`
			// (preview) `'Bob'`
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5Gconst replScopedVariable = 'Bob'
\x1b[90mundefined\x1b[39m
\x1b[1G\x1b[0J--> \x1b[5GreplScopedVariable
\x1b[90m'Bob'\x1b[39m\x1b[23G\x1b[1A
`),
		},
		{
			name:  "global variable",
			input: "obj",
			// (prompt) `obj`
			// (preview) `{ abc: 'def' }`
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5Gobj
\x1b[90m{ abc: 'def' }\x1b[39m\x1b[8G\x1b[1A
`),
		},
		{
			name:  "new *",
			input: "new MyClass()",
			// (prompt) `new MyClass()`
			// (preview) `MyClass {}`
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5Gnew MyClass()
\x1b[90mMyClass {}\x1b[39m\x1b[18G\x1b[1A
`),
		},
		{
			name: "global variable + module variable + repl variable",
			input: trimNewlines(`
const c = 23
a + b + c
`),
			// (prompt)  `const c = 23;`
			// (result)  `undefined`
			// (prompt)  `a + b + c`
			// (preview) `60`
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5Gconst c = 23
\x1b[90mundefined\x1b[39m
\x1b[1G\x1b[0J--> \x1b[5Ga + b + c
\x1b[90m60\x1b[39m\x1b[14G\x1b[1A
`),
		},
		{
			name:  "property of global object",
			input: "obj.abc",
			// (prompt)  `obj.abc`
			// (preview) `'def'`
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5Gobj.abc
\x1b[90m'def'\x1b[39m\x1b[12G\x1b[1A
`),
		},
		{
			name:  "assignment to global variable",
			input: "b = 3",

			// (prompt) `b = 3`
			expected: "\x1b[1G\x1b[0J--> \x1b[5Gb = 3",
		},
		{
			name:  "assignment to module variable",
			input: "obj = 23",

			// (prompt) `obj = 23`
			expected: "\x1b[1G\x1b[0J--> \x1b[5Gobj = 23",
		},
		{
			name:  "assignment to constant module variable",
			input: "obj = 23",

			// (prompt) `obj = 23`
			expected: "\x1b[1G\x1b[0J--> \x1b[5Gobj = 23",
		},
		{
			name: "displays the cached value of a local variable",
			input: trimNewlines(`
setValue(2);
value
`),
			// (prompt) `setValue(2);`
			// (result) `undefined`
			// (prompt) `value`
			// (preview) `1`
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5GsetValue(2);
\x1b[90mundefined\x1b[39m
\x1b[1G\x1b[0J--> \x1b[5Gvalue
\x1b[90m1\x1b[39m\x1b[10G\x1b[1A
`),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.run(t, info)
		})
	}
}

func testEvalResults(t *testing.T, info testInfo) {
	t.Parallel()

	for _, tc := range [...]testCase{
		{
			name:  "number",
			input: "3\n",
			// (prompt) `3`
			// (result) `3`
			// (prompt)
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5G3
\x1b[33m3\x1b[39m
\x1b[1G\x1b[0J--> \x1b[5G
`),
		},
		{
			name:  "number + number",
			input: "5 + 10\n",
			// (prompt) `5 + 10`
			// (result) `15`
			// (prompt)
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5G5 + 10
\x1b[33m15\x1b[39m
\x1b[1G\x1b[0J--> \x1b[5G
`),
		},
		{
			name:  "module variable + number",
			input: "a + 3\n",
			// (prompt) `a + 3`
			// (result) `8`
			// (prompt)
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5Ga + 3
\x1b[33m8\x1b[39m
\x1b[1G\x1b[0J--> \x1b[5G
`),
		},
		{
			name:  "non-existant variable",
			input: "notAVariable\n",
			// (prompt) `notAVariable`
			// (result) (error)
			// (prompt)
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5GnotAVariable
\x1b[0m\x1b[31mReferenceError: notAVariable is not defined\x1b[0m
\x1b[1G\x1b[0J--> \x1b[5G
`),
		},
		{
			name: "repl-scoped variable",
			input: trimNewlines(`
const replScopedVariable = 'Bob'
replScopedVariable

`),
			// (prompt) `const replScopedVariable = 'Bob'`
			// (result) `undefined`
			// (prompt) `replScopedVariable`
			// (result) `'Bob'`
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5Gconst replScopedVariable = 'Bob'
\x1b[90mundefined\x1b[39m
\x1b[1G\x1b[0J--> \x1b[5GreplScopedVariable
\x1b[32m'Bob'\x1b[39m
\x1b[1G\x1b[0J--> \x1b[5G
`),
		},
		{
			name:  "global variable",
			input: "obj\n",
			// (prompt) `obj`
			// (result) `{ abc: 'def' }`
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5Gobj
{ abc: \x1b[32m'def'\x1b[39m }
\x1b[1G\x1b[0J--> \x1b[5G
`),
		},
		{
			name:  "new *",
			input: "new MyClass()\n",
			// (prompt) `new MyClass()`
			// (result) `MyClass {}`
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5Gnew MyClass()
MyClass {}
\x1b[1G\x1b[0J--> \x1b[5G
`),
		},
		{
			name: "global variable + module variable + repl variable",
			input: trimNewlines(`
const c = 23
a + b + c

`),
			// (prompt) `const c = 23;`
			// (result) `undefined`
			// (prompt) `a + b + c`
			// (result) `60`
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5Gconst c = 23
\x1b[90mundefined\x1b[39m
\x1b[1G\x1b[0J--> \x1b[5Ga + b + c
\x1b[33m60\x1b[39m
\x1b[1G\x1b[0J--> \x1b[5G
`),
		},
		{
			name:  "property of global object",
			input: "obj.abc\n",
			// (prompt) `obj.abc`
			// (result) `'def'`
			// (prompt)
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5Gobj.abc
\x1b[32m'def'\x1b[39m
\x1b[1G\x1b[0J--> \x1b[5G
`),
		},
		{
			name:  "assignment to a local variable",
			input: "b = 3\n",
			// (prompt) `b = 3`
			// (result) `3`
			// (prompt)
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5Gb = 3
\x1b[33m3\x1b[39m
\x1b[1G\x1b[0J--> \x1b[5G
`),
		},
		{
			name:  "assignment to module variable",
			input: "obj = 23\n",

			// (prompt) `obj = 23`
			// (result) `23`
			// (prompt)
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5Gobj = 23
\x1b[33m23\x1b[39m
\x1b[1G\x1b[0J--> \x1b[5G
`),
		},
		{
			name:  "direct console.trace",
			input: "console.trace('Hello')\n",

			// (prompt) `obj = 23`
			// (trace) `Trace: Hello`
			// (result) `23`
			// (prompt)
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5Gconsole.trace('Hello')
Trace: Hello
\x1b[90mundefined\x1b[39m
\x1b[1G\x1b[0J--> \x1b[5G
`),
		},

		{
			name: "gets the latest value of a variable",
			input: trimNewlines(`
setValue(2);
value

`),
			// (prompt) `setValue(2);`
			// (result) `undefined`
			// (prompt) `value`
			// (result) `2`
			// (prompt)
			expected: unescapeExpectedOutput(`
\x1b[1G\x1b[0J--> \x1b[5GsetValue(2);
\x1b[90mundefined\x1b[39m
\x1b[1G\x1b[0J--> \x1b[5Gvalue
\x1b[33m2\x1b[39m
\x1b[1G\x1b[0J--> \x1b[5G
`),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.run(t, info)
		})

	}
}

func TestREPL(t *testing.T) {
	const code = `
const str = 'Hello, World!';
const a = 5;
global.b = 32;
let obj = { abc: 'def' };
class MyClass {}

let value = 1;

function setValue(newValue) {
	value = 2;
}
`

	for _, info := range [...]testInfo{
		{
			name: "string",
			args: getCodeArgs(true, true, code),
		},
		{
			name:   "file",
			args:   getFileArgs(t, true, true, code),
			prefix: "\x1b[0m\x1b[90mHint: hit control+c anytime to enter REPL.\x1b[0m\n",
		},
	} {
		for _, test := range [...]struct {
			run  func(t *testing.T, info testInfo)
			name string
		}{
			{
				name: "result preview",
				run:  testPreview,
			},
			{
				name: "line result",
				run:  testEvalResults,
			},
		} {
			t.Run(test.name+"/"+info.name, func(t *testing.T) {
				test.run(t, info)
			})
		}
	}

}
