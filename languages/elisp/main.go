package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/replit/prybar/utils"
)

func Execute(config *utils.Config) {
	args := []string{}

	if _, err := os.Stat("Cask"); err == nil {
		if _, err := exec.LookPath("cask"); err == nil {
			args = append(args, "cask", "exec")
		}
	}

	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	runDir := filepath.Dir(execPath)
	replPath := filepath.Join(runDir, "prybar_assets", "elisp", "repl.el")

	args = append(
		args, "emacs", "-nw", "-Q", "--load", replPath,
		"--eval", "(prybar-repl)",
	)

	os.Setenv("PRYBAR_EVAL", config.Exp)
	os.Setenv("PRYBAR_EXEC", config.Code)
	os.Setenv("PRYBAR_PS1", config.Ps1)

	if config.Quiet {
		os.Setenv("PRYBAR_QUIET", "1")
	} else {
		os.Setenv("PRYBAR_QUIET", "")
	}

	// We only support one file, despite the fact that this
	// variable is a list.
	if len(config.Args) >= 1 {
		os.Setenv("PRYBAR_FILE", config.Args[0])
	} else {
		os.Setenv("PRYBAR_FILE", "")
	}

	if !(config.Interactive || config.OurInteractive) {
		fmt.Fprintln(os.Stderr, "prybar-elisp: warn: non-interactive mode not implemented")
	}

	if config.Ps2 != "... " {
		fmt.Fprintln(os.Stderr, "prybar-elisp: warn: ps2 not implemented")
	}

	path, err := exec.LookPath(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: command not found\n", args[0])
	}

	err = syscall.Exec(path, args, os.Environ())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to exec %s: %s\n", path, err)
		os.Exit(1)
	}

}
