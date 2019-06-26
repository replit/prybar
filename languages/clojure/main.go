package main

import (
	"os"
	"os/exec"
	"syscall"
	"github.com/replit/prybar/utils"
)

func Execute(config *utils.Config) {
	path, err := exec.LookPath("clj")

	if err != nil {
		panic(err)
	}

	env := os.Environ()
	args := []string{"clj"}

	if config.Code != "" {
		args = append(args, "--eval", config.Code)
	}

	if config.Exp != "" {
		args = append(args, "--eval", config.Exp)
	}

	if (config.Interactive || config.OurInteractive) {
		args = append(args, "--repl", config.Exp)
	}

	if config.Ps1 != "" {
		env = append(env, "PRYBAR_PS1="+config.Ps1)
	}

	if config.Ps2 != "" {
		env = append(env, "PRYBAR_PS2="+config.Ps2)
	}

	if config.Quiet {
		// no-op:
		// not supported by `clojure`.
	}

	if config.Args != nil && len(config.Args) > 0 {
		args = append(args, config.Args...)
	} else if config.Exp == "" && config.Code == "" {
		args = append(args, "--eval", "")
	}

	syscall.Exec(path, args, env)
}

