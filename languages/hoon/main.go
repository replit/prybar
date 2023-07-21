package main

// USING_CGO

import (
	"github.com/chzyer/readline"

	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Hoon struct{}

func (p Hoon) Open() {
}

func (p Hoon) SetPrompts(ps1, ps2 string) {
	p.ps1 = ps1
	p.ps2 = ps2
}

func GetPath() string {
	pat, err := exec.LookPath("urbit")
	if err != nil {
		panic(err)
	}
	return pat
}

func (p Hoon) Version() string {
	urb := exec.Command(GetPath(), "--version")
	versions, err := urb.Output()
	if err != nil {
		panic(err)
	}

	return string(versions)
}

func (p Hoon) EvalFile(file string, args []string) int {
	fileContents, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	_, _, err = RunCommand(string(fileContents))
	if err != nil {
		panic(err)
	}

	return 0
}

func RunCommand(command string) (string, bool, error) {
	urb := exec.Command(GetPath(), "eval")
	stdin, err := urb.StdinPipe()
	if err != nil {
		return "", false, err
	}
	io.WriteString(stdin, command+"\n")
	stdin.Close()
	out, err := urb.Output()
	if err != nil {
		return "", false, err
	}

	stringOut := string(out)

	lines := strings.Split(stringOut, "\n")

	needsMoreInput := true
	if len(lines) > 1 && len(lines[0]) == 21 {
		needsMoreInput = needsMoreInput && (lines[0][11:16] == "/eval")

		rxp := regexp.MustCompile(`\{(\d+) (\d+)\}`)
		lineMatches := rxp.FindSubmatch([]byte(lines[1]))
		if len(lineMatches) == 3 {
			l, _ := strconv.Atoi(string(lineMatches[1]))
			c, _ := strconv.Atoi(string(lineMatches[2]))

			realL := len(strings.Split(command, "\n"))

			needsMoreInput = needsMoreInput && (realL+1 == l) && (c == 1)

		} else {
			needsMoreInput = false
		}

	} else {
		needsMoreInput = false
	}

	return stringOut, needsMoreInput, nil

}

func (p Hoon) EvalExpression(code string) string {
	out, _, err := RunCommand(code)

	if err != nil {
		panic(err)
	}

	return string(out)
}

func (p Hoon) REPL() {
	for {
		line, err := readline.Line("--> ")
		if err != nil {
			break
		}
		readline.AddHistory(line)
		out, needMoreInput, err := RunCommand(line)

		for needMoreInput {
			newLine, err := readline.Line("... ")
			if err != nil {
				break
			}
			readline.AddHistory(newLine)
			line = line + "\n" + newLine
			out, needMoreInput, err = RunCommand(line)
		}

		if err != nil {
			panic(err)
		}

		strOut := string(out)
		fmt.Println(strOut)
	}
}

func (p Hoon) Close() {
}

var Instance = Hoon{}
