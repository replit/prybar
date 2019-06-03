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

	flag.StringVar(&config.Code, "c", "", "execute without printing result")
	flag.StringVar(&config.Exp, "e", "", "evaluate and print result")

	flag.StringVar(&config.Ps1, "ps1", "--> ", "repl prompt")
	flag.StringVar(&config.Ps2, "ps2", "... ", "repl continuation prompt")

	flag.BoolVar(&config.Interactive, "i", false, "interactive (use language repl)")
	flag.BoolVar(&config.OurInteractive, "I", false, "interactive (use readline repl)")
	flag.BoolVar(&config.Quiet, "q", false, "don't print language version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [FLAGS] [FILENAME]...\n", os.Args[0])
		flag.PrintDefaults()
	}

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
