// OCaml main
package main

//go:generate ../../scripts/gofiles.sh generated_files.go

import (
	"github.com/replit/prybar/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
)

func findHelper(path string) string {
	bytes, err := File(path + ".ml")
	if bytes != nil {
		f, err := ioutil.TempFile("", "prybar-ocaml-*.ml")
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
	panic("File not found: " + path + ".ml")
}

func Execute(config *utils.Config) {
	path, err := exec.LookPath("ocaml")

	if err != nil {
		panic(err)
	}

	env := os.Environ()
	args := []string{"ocaml", findHelper("repl"), "-s", "ml"}

	if config.Quiet {
		args = append(args, "-q")
	}

	if config.Code != "" {
		args = append(args, "-c", config.Code)
	} else if config.Exp != "" {
		args = append(args, "-e", config.Exp)
	} else if (config.Interactive || config.OurInteractive) {
		args = append(args, "-i")
	}

	if config.Ps1 != "" {
		env = append(env, "PRYBAR_PS1="+config.Ps1)
	}

	if config.Ps2 != "" {
		env = append(env, "PRYBAR_PS2="+config.Ps2)
	}

	if config.Args != nil && len(config.Args) > 0 {
		args = append(args, config.Args...)
	}

	syscall.Exec(path, args, env)
}
