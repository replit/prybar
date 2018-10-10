package main

import (
	"fmt"
	"os"
	"plugin"
	"flag"
)

type Language interface {
	Open()
	Close()
	Eval(string)
	Version() string
	REPL()
}


func main() {
	var language string
	var interactive bool
	var code string
	var quiet bool


	flag.StringVar(&language, "l", "python2", "langauge")
	flag.StringVar(&code, "c", "", "code to run")
	flag.BoolVar(&interactive, "i", false, "interactive")
	flag.BoolVar(&quiet, "q", false, "quiet")
	flag.Parse()


	plug, err := plugin.Open("./plugins/" + language + ".so")

	symGreeter, err := plug.Lookup("Instance")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var lang Language
	lang, ok := symGreeter.(Language)
	if !ok {
		fmt.Println("unexpected type from module symbol")
		os.Exit(1)
	}

	// 4. use the module
	lang.Open()
	if !quiet {
		fmt.Println(lang.Version())
	}
	if code != "" {
		lang.Eval(code)
	}
	if interactive {
		lang.REPL()
	}
	lang.Close()
	
}

var instance Language