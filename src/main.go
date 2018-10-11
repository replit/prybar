package main

import (
	"fmt"
	"flag"
)




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

	args := flag.Args()

	

	// 4. use the module
	lang := GetLanguage(language)
	if !quiet {
		fmt.Println(lang.Version())
	}
	if code != "" {
		lang.Eval(code)
	} else if len(args) > 0 {
		lang.EvalFile(args[0], args[1:])
	}
	if interactive {
		lang.REPL()
	}

	
}

