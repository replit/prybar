package main

//go:generate ../../scripts/gofiles.sh generated_files.go

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"github.com/replit/prybar/utils"
)

func findHelper(path string) string {
    fileName := path + ".clj"
	bytes, err := File(fileName)
	if bytes != nil {
		f, err := ioutil.TempFile("", fileName)
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
    panic("File not found: " + fileName)
}

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

	if config.Ps1 != "" {
		args = append(args, "-J-DPRYBAR_PS1=" + config.Ps1)
	}

	if config.Ps2 != "" {
		args = append(args, "-J-DPRYBAR_PS2=" + config.Ps2)
	}

	if config.Quiet {
		args = append(args, "-J-DPRYBAR_QUIET=true")
	}

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
			// "--eval" prints non-nil results only.
			effect := fmt.Sprintf("(do %s nil)", config.Code)
			args = append(args, "--eval", effect)
		}

		if config.Exp != "" {
			args = append(args, "--eval", config.Exp)
		}

		if config.Interactive || config.OurInteractive {
			args = append(args, findHelper("prybar_repl"))
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

	syscall.Exec(path, args, env)
}

