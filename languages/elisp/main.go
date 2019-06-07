package main

import (
	"fmt"
	"github.com/replit/prybar/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func Execute(config *utils.Config) {

	emacs, err := exec.LookPath("emacs")

	if err != nil {
		panic(err)
	}

	execPath, err := os.Executable()

	if err != nil {
		panic(err)
	}

	runDir := filepath.Dir(execPath)
	replPath := filepath.Join(runDir, "prybar_assets", "elisp", "repl.el")

	os.Setenv("PRYBAR_EVAL", config.Exp)
	os.Setenv("PRYBAR_EXEC", config.Code)
	os.Setenv("PRYBAR_PS1", config.Ps1)

	if config.Quiet {
		os.Setenv("PRYBAR_QUIET", "1")
	} else {
		os.Setenv("PRYBAR_QUIET", "")
	}

	if config.Args != nil {
		os.Setenv("PRYBAR_FILES", strings.Join(config.Args, "\000"))
	} else {
		os.Setenv("PRYBAR_FILES", "")
	}

	if !(config.Interactive || config.OurInteractive) {
		fmt.Fprintln(os.Stderr, "prybar-elisp: warn: non-interactive mode not implemented");
	}

	if config.Ps2 != "... " {
		fmt.Fprintln(os.Stderr, "prybar-elisp: warn: ps2 not implemented");
	}

	args := []string{"emacs", "-Q", "--load", replPath, "--eval", "(prybar-repl)"}
	err = syscall.Exec(emacs, args, os.Environ())

	if err != nil {
		panic(err)
	}

}
