package main

import (
	"os"
	"os/exec"
	"syscall"
)

func Execute(config *utils.Config) {

	os.Setenv("PRYBAR_EVAL", config.Exp)
	os.Setenv("PRYBAR_EXEC", config.Code)
	os.Setenv("PRYBAR_PS1", config.Ps1)
	os.Setenv("PRYBAR_PS2", config.Ps2)
	os.Setenv("PRYBAR_QUIET", config.Quiet)

	if config.Interactive || config.OurInteractive {
		os.Setenv("PRYBAR_INTERACTIVE", "1")
	} else {
		os.Setenv("PRYBAR_INTERACTIVE", "")
	}

	args := []string{"emacs", "-Q", "--script", "repl.el"}
	syscall.Exec("emacs", args, os.Environ())

}
