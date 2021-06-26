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
	bytes, err := File(path)

	if err != nil {
		panic(err)
	}

	if bytes == nil {
		panic("File not found at path: " + path)
	}

	f, err := ioutil.TempFile("", path)

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

func Execute(config *utils.Config) {
	cljPath, err := exec.LookPath("clj")

	if err != nil {
		panic(err)
	}

	env := os.Environ()
	args := []string{
		"clj",
		"-Sdeps",
		"{:deps {org.clojure/tools.namespace {:mvn/version \"1.0.0\"}} :paths [\"src\" \".\"]}"}

	hasAction := config.Code != "" || config.Exp != "" ||
		config.Interactive || config.OurInteractive
	filePath := ""

	if config.Ps1 != "" {
		args = append(args, fmt.Sprintf("-J-DPRYBAR_PS1=%s", config.Ps1))
	}

	if config.Quiet {
		args = append(args, "-J-DPRYBAR_QUIET=true")
	}

	if !hasAction {
		// An empty eval-opt suppresses the start of a REPL.
		args = append(args, "--eval", "")
		args = append(args, config.Args...)
		syscall.Exec(cljPath, args, env)

		return
	}

	if config.Args != nil && len(config.Args) > 0 {
		if _, err := os.Stat(config.Args[0]); os.IsNotExist(err) {
			fmt.Println("No such file:", config.Args[0])
			os.Exit(2)
		}

		filePath = config.Args[0]
	}

	if config.Code != "" {
		// "--eval" prints non-nil results only.
		effect := fmt.Sprintf("(do %s nil)", config.Code)
		args = append(args, "--eval", effect)
	}

	if config.Exp != "" {
		args = append(args, "--eval", config.Exp)
	}

	if !(config.Interactive || config.OurInteractive) {
		args = append(args, config.Args...)
		syscall.Exec(cljPath, args, env)

		return
	}

	// Starting a REPL, pass the file along, if exists, as an init-opt.
	if filePath != "" {
		args = append(args, "--init", filePath)
	}

	args = append(args, findHelper("prybar_clojure_repl.clj"))

	if filePath != "" {
		args = append(args, config.Args[1:]...)
	} else {
		args = append(args, config.Args...)
	}

	syscall.Exec(cljPath, args, env)
}

