package main

import (
	"fmt"
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

	hasOption := config.Code != "" || config.Exp != "" ||
		config.Interactive || config.OurInteractive
	hasFile := false

	if hasOption {
		if config.Args != nil && len(config.Args) > 0 {
			if _, err := os.Stat(config.Args[0]); os.IsNotExist(err) {
				fmt.Println("No such file:", config.Args[0])
				os.Exit(2)
			}

			hasFile = true
		}

		if hasFile {
			args = append(args, "--init", config.Args[0])
		}

		if config.Code != "" {
			args = append(args, "--eval", config.Code)
		}

		if config.Exp != "" {
			args = append(args, "--eval", config.Exp)
		}

		if config.Interactive || config.OurInteractive {
			// "The appearance of any eval option before running a repl
			// suppresses the usual greeting message: \"Clojure ~(clojure-version)\"."
			// (https://github.com/clojure/clojure/blob/653b8465845a78ef7543e0a250078eea2d56b659/src/clj/clojure/main.clj#L644)
			if config.Quiet {
				args = append(args, "--eval", "")
			}

			args = append(args, "--repl")
		}

		if hasFile {
			args = append(args, config.Args[1:]...)
		} else {
			args = append(args, config.Args...)
		}
	} else {
		args = append(args, "--eval", "")
		args = append(args, config.Args...)
	}

	if config.Ps1 != "" {
		env = append(env, "PRYBAR_PS1="+config.Ps1)
	}

	if config.Ps2 != "" {
		env = append(env, "PRYBAR_PS2="+config.Ps2)
	}

	syscall.Exec(path, args, env)
}

