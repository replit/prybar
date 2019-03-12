package main

import (
	"fmt"

	"github.com/cosmos72/gomacro/base"
	"github.com/cosmos72/gomacro/fast"
	"github.com/replit/prybar/utils"
)

func Execute(config *utils.Config) {

	interp := fast.New()
	interp.Comp.Prompt = config.Ps1
	interp.Comp.Globals.Options = base.OptDebugger | base.OptCtrlCEnterDebugger | base.OptKeepUntyped | base.OptTrapPanic | base.OptShowPrompt | base.OptShowEval | base.OptShowEvalType
	interp.Comp.Globals.Options |= base.OptShowPrompt

	if config.Code != "" {

		// run a file
		_, err := interp.EvalFile(config.Code)
		if err != nil {

			fmt.Println("error:", err)

		}

	} else if config.Exp != "" {

		// run an expression
		vals, types := interp.Eval(config.Exp)
		for i, val := range vals {

			fmt.Printf("%v // %v\n", val, types[i])

		}

	} else {

		// run the repl
		interp.ReplStdin()

	}

}
