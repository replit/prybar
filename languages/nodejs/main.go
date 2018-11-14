package main

//go:generate ../../scripts/gofiles.sh generated_files.go

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

	env := os.Environ()
	args := []string{"node"}

	if !config.Quiet {
		args = append(args, "-r", findHelper("version"))

	}

	if config.Code != "" {
		args = append(args, "-e", config.Code)
	}

	if config.Exp != "" {
		args = append(args, "-p", config.Exp)
	}

	if config.Interactive {
		args = append(args, "-r", findHelper("repl"))
	}

	if config.Ps1 != "" {
		env = append(env, "NODE_PROMPT="+config.Ps1)
	}

	if config.Args != nil && len(config.Args) > 0 {
		args = append(args, config.Args...)
	} else if config.Exp == "" && config.Code == "" {
		args = append(args, "-e", "")
	}

	syscall.Exec(path, args, env)
}
