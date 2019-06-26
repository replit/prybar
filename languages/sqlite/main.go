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

// constructConfigFile generates commands to configure the sqlite CLI.
// It writes them to a temporary file and returns its pathname.
func constructConfigFile(config *utils.Config) string {
	f, err := ioutil.TempFile("", "sqlite-config")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// main and continuation prompts
	// TODO: this probably doesn't handle quotation marks properly
	f.WriteString(fmt.Sprintf(".prompt '%s' '%s'\n", config.Ps1, config.Ps2))

	// execute file, if specified
	if len(config.Args) == 1 {
		fileToRun := config.Args[0]
		f.WriteString(fmt.Sprintf(".read %s\n", fileToRun))
	}

	return f.Name()
}

func evalAndPrint(sqlite string, config *utils.Config) {
	args := []string{"sqlite3"}
	env := os.Environ()

	if len(config.Code) > 0 {
		f, err := ioutil.TempFile("", "sqlite-config")
		if err != nil {
			panic(err)
		}
		f.WriteString(".output /dev/null\n")
		f.Close()
		args = append(args, []string{"-init", f.Name(), ":memory:", config.Code}...)

		// add LD_PRELOAD lib to environment to suppress initialization output
		execPath, err := os.Executable()
		if err != nil {
			panic(err)
		}
		runDir := filepath.Dir(execPath)
		libPath := filepath.Join(runDir, "prybar_assets", "sqlite", "patch.so")
		env = append(env, []string{"LD_PRELOAD="+libPath, "PRYBAR_QUIET=1"}...)
	} else {
		args = append(args, []string{":memory:", config.Exp}...)
	}

	err := syscall.Exec(sqlite, args, env)
	if err != nil {
		panic(err)
	}
}

func interactive(sqlite string, config *utils.Config) {
	configFile := constructConfigFile(config)
	args := []string{"sqlite3", "-init", configFile}

	env := os.Environ()
	if config.Quiet {
		env = append(env, "PRYBAR_QUIET=1")
	}

	// add LD_PRELOAD lib to environment
	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	runDir := filepath.Dir(execPath)
	libPath := filepath.Join(runDir, "prybar_assets", "sqlite", "patch.so")
	env = append(env, "LD_PRELOAD="+libPath)

	err = syscall.Exec(sqlite, args, env)
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

	sqlite, err := exec.LookPath("sqlite3")
	if err != nil {
		panic(err)
	}

	if len(config.Code) > 0 || len(config.Exp) > 0 {
		evalAndPrint(sqlite, config)
	} else {
		interactive(sqlite, config)
	}

}
