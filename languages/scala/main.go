package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"github.com/replit/prybar/utils"
)


func Execute(config *utils.Config) {
	path, err := exec.LookPath("scala")

	if err != nil {
		panic(err)
	}

	env := os.Environ()
	args := []string{"scala"}

	if config.Quiet {
		args = append(args, "-Dscala.shell.welcome=")
	}

	if config.Code != "" {
		args = append(args, "-e", config.Code)
	}

	if config.Exp != "" {
		printStatement := fmt.Sprintf("print(%s)", config.Exp)
		args = append(args, "-e", printStatement)
	}

	if config.Ps1 != "" {
		args = append(args, "-Dscala.shell.prompt=%n" + config.Ps1)
	}

	if config.Args != nil && len(config.Args) > 0 {
		args = append(args, "-i")
		args = append(args, config.Args...)
	}

	syscall.Exec(path, args, env)
}
