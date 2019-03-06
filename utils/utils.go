package utils

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

type Config struct {
	Code                        string
	Quiet                       bool
	Exp                         string
	Ps1, Ps2                    string
	Args                        []string
	Interactive, OurInteractive bool
}

func ParseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.Code, "c", "", "code to run")
	flag.StringVar(&config.Exp, "e", "", "expression to print")

	flag.StringVar(&config.Ps1, "ps1", "--> ", "PS1")
	flag.StringVar(&config.Ps2, "ps2", "... ", "PS2")

	flag.BoolVar(&config.Interactive, "i", false, "interactive")
	flag.BoolVar(&config.OurInteractive, "I", false, "like -i, but never use language REPL")
	flag.BoolVar(&config.Quiet, "q", false, "quiet")

	flag.Parse()

	config.Args = flag.Args()
	return config
}

func DoCli(p PluginBase) {

	config := ParseFlags()

	runtime.LockOSThread()
	p.Open()
	lang := &Language{ptr: p, ps1: config.Ps1}

	if !config.Quiet {
		fmt.Println(lang.Version())
	}

	if config.Code != "" {
		lang.Eval(config.Code)
	}
	if config.Exp != "" {
		lang.EvalAndTryToPrint(config.Exp)
	}
	if len(config.Args) > 0 {
		if _, err := os.Stat(config.Args[0]); os.IsNotExist(err) {
			fmt.Println("No such file:", config.Args[0])
			os.Exit(2)
		} else {
			lang.EvalFile(config.Args[0], config.Args[1:])
		}
	}

	if config.Interactive {
		lang.SetPrompts(config.Ps1, config.Ps2)
		lang.REPL()
	} else if config.OurInteractive {
		lang.InternalREPL()
	}
}
