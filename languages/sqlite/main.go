// Take a Prybar config and exec the SQLite CLI with it.
//
// Note: the prompt config (ps1 and ps2) must be SQL strings, i.e. single quotation
// marks must be properly escaped.
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/replit/prybar/utils"
)

var Instance = &SQLite{}

type SQLite struct{}

// writeConfig writes the provided slice of strings out to temporary file and
// returns its pathname.
func writeConfig(lines []string) string {
	f, err := ioutil.TempFile("", "sqlite-config")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for _, line := range lines {
		f.WriteString(line)
	}

	return f.Name()
}

// preloadQuietLib adds our LD_PRELOAD lib to environment and enables it
func preloadQuietLib(env []string) []string {
	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	runDir := filepath.Dir(execPath)
	libPath := filepath.Join(runDir, "prybar_assets", "sqlite", "patch.so")
	return append(env, []string{"LD_PRELOAD=" + libPath, "PRYBAR_QUIET=1"}...)
}

func eval(sqlite string, config *utils.Config) {
	args := []string{"sqlite3"}
	env := os.Environ()

	if len(config.Code) > 0 {
		// We have to execute provided code without showing output.
		// So, let's have sqlite run a command that will output to /dev/null.
		confLines := []string{".output /dev/null\n"}
		args = append(args, []string{"-init", writeConfig(confLines), ":memory:", config.Code}...)

		// add LD_PRELOAD lib to environment to suppress initialization output
		env = preloadQuietLib(env)
	} else {
		// we can just run the code without a config file since we want its output
		args = append(args, []string{":memory:", config.Exp}...)
	}

	err := syscall.Exec(sqlite, args, env)
	if err != nil {
		panic(err)
	}
}

func interactive(sqlite string, config *utils.Config) {
	// main and continuation prompts
	confLines := []string{
		fmt.Sprintf(".prompt '%s' '%s'\n", config.Ps1, config.Ps2),
	}

	// execute file, if specified
	if len(config.Args) == 1 {
		fileToRun := config.Args[0]
		confLines = append(confLines, fmt.Sprintf(".read %s\n", fileToRun))
	}

	args := []string{"sqlite3", "-init", writeConfig(confLines)}

	env := os.Environ()
	if config.Quiet {
		env = preloadQuietLib(env)
	}

	err := syscall.Exec(sqlite, args, env)
	if err != nil {
		panic(err)
	}
}

func Execute(config *utils.Config) {
	// sanity check
	if len(config.Args) > 1 {
		fmt.Fprint(os.Stderr, "too many arguments\n")
		os.Exit(1)
	}

	if config.OurInteractive {
		fmt.Fprint(os.Stderr, "not supported\n")
		os.Exit(1)
	}

	// make sure we have the sqlite CLI
	sqlite, err := exec.LookPath("sqlite3")
	if err != nil {
		panic(err)
	}

	if len(config.Code) > 0 || len(config.Exp) > 0 {
		// we're executing provided code
		eval(sqlite, config)
	} else {
		// running sqlite CLI interactively
		interactive(sqlite, config)
	}

}
