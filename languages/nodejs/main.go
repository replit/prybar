package main

//go:generate bash ../../scripts/gofiles.sh generated_files.go

import (
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	"github.com/replit/prybar/utils"
)

func findHelper(path string) string {
	bytes, err := File(path + ".js")
	if bytes != nil {
		f, err := ioutil.TempFile("", "prybar-nodejs-*.js")
		if err != nil {
			panic(err)
		}

		if _, err = f.Write(bytes); err != nil {
			panic(err)
		}
		if err = f.Close(); err != nil {
			panic(err)
		}
		return f.Name()
	}
	if err != nil {
		panic(err)
	}
	panic("File not found")
}

func Execute(config *utils.Config) {
	path, err := exec.LookPath("node")

	if err != nil {
		panic(err)
	}

	args := []string{"node", findHelper("repl")}

	if config.Quiet {
		os.Setenv("PRYBAR_QUIET", "1")
	}

	os.Setenv("PRYBAR_CODE", config.Code)
	os.Setenv("PRYBAR_EXP", config.Exp)
	os.Setenv("PRYBAR_PS1", config.Ps1)

	if config.Interactive {
		os.Setenv("PRYBAR_INTERACTIVE", "1")
	}

	// We only support one file, despite the fact that this
	// variable is a list.
	if len(config.Args) >= 1 {
		os.Setenv("PRYBAR_FILE", config.Args[0])
	} else {
		os.Setenv("PRYBAR_FILE", "")
	}

	syscall.Exec(path, args, os.Environ())
}
