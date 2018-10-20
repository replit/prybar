package utils

import (
	"flag"
	"fmt"
	"io"
	"os"
	"syscall"
)

type Red struct {
	parent io.Writer
}

func (r *Red) Write(p []byte) (n int, err error) {
	r.parent.Write(([]byte("\033[31m")))
	n, err = r.parent.Write(p)
	r.parent.Write(([]byte("\033[0m")))
	return n, err
}

func DoCli(p PluginBase) {
	var interactive, ourInteractive bool
	var code string
	var quiet bool
	var exp string
	var colorizeStderr bool
	var ps1, ps2 string

	flag.StringVar(&code, "c", "", "code to run")
	flag.StringVar(&exp, "e", "", "expression to print")

	flag.StringVar(&ps1, "ps1", "--> ", "PS1")
	flag.StringVar(&ps2, "ps2", "... ", "PS2")

	flag.BoolVar(&interactive, "i", false, "interactive")
	flag.BoolVar(&ourInteractive, "I", false, "like -i, but never use language REPL")
	flag.BoolVar(&quiet, "q", false, "quiet")

	flag.BoolVar(&colorizeStderr, "R", false, "color standard error red")

	flag.Parse()

	args := flag.Args()

	p.Open()
	lang := &Language{ptr: p, ps1: ps1}

	if !quiet {
		fmt.Println(lang.Version())
	}

	if colorizeStderr {
		var pipes [2]int
		newStderr, err := syscall.Dup(2)
		if err != nil {
			panic(err)
		}
		syscall.Pipe(pipes[:])
		syscall.Dup2(pipes[1], 2)
		syscall.Close(pipes[1])
		o := &Red{parent: os.NewFile(uintptr(newStderr), "o")}
		i := os.NewFile(uintptr(pipes[0]), "i")
		go io.Copy(o, i)
	}

	if code != "" {
		lang.Eval(code)
	}
	if exp != "" {
		lang.EvalAndTryToPrint(exp)
	}
	if len(args) > 0 {
		if _, err := os.Stat(args[0]); os.IsNotExist(err) {
			fmt.Println("No such file:", args[0])
			os.Exit(2)
		} else {
			lang.EvalFile(args[0], args[1:])
		}
	}

	if interactive {
		lang.SetPrompts(ps1, ps2)
		lang.REPL()
	} else if ourInteractive {
		//LinenoiseSetCompleter(func(s string) []string {
		//	return []string{s + "A", s + "B", s + "B"}
		//})
		lang.InternalREPL()
	}

}
