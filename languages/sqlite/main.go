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

	return f.Name()
}

func Execute(config *utils.Config) {
	// sanity check
	if len(config.Args) > 1 {
		panic("too many arguments")
	}

	sqlite, err := exec.LookPath("sqlite3")
	if err != nil {
		panic(err)
	}

	configFile := constructConfigFile(config)
	args := []string{"sqlite3", "-init", configFile}

	env := os.Environ()
	if config.Quiet {
		env = append(env, "PRYBAR_QUIET=true")
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

	// execute file, if specified
	// if len(config.Args) == 1 {
	// 	fileToRun := config.Args[0]
	// 	ptty.WriteString(fmt.Sprintf(".read %s\n", fileToRun))
	// }

}
