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
			args = append(args, "--i", config.Args[0])
		}

		if config.Exp != "" {
			args = append(args, "--eval", config.Exp)
		}

		if hasFile {
			args = append(args, config.Args[1:]...)
		} else {
			args = append(args, config.Args...)
		}
	} else {
		args = append(args, config.Args...)
	}

	syscall.Exec(path, args, env)
}

