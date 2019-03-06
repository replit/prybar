package utils

import (
	"fmt"
	"io/ioutil"

	"github.com/chzyer/readline"
)

type PluginBase interface {
	Open()
	Close()
	Version() string
}

type PluginEval interface {
	PluginBase
	Eval(code string)
}

type PluginEvalExpression interface {
	PluginBase
	EvalExpression(code string) string
}

type PluginEvalFile interface {
	PluginBase
	EvalFile(file string, args []string)
}

type PluginREPL interface {
	PluginBase
	REPL()
}

type PluginREPLLikeEval interface {
	PluginBase
	REPLLikeEval(code string)
}

type PluginSetPrompts interface {
	PluginBase
	SetPrompts(ps1, ps2 string)
}

type Language struct {
	ptr PluginBase
	ps1 string
}

func (lang Language) Version() string {
	return lang.ptr.Version()
}

func (lang Language) Eval(code string) {
	lang.ptr.(PluginEval).Eval(code)
}

func (lang Language) EvalAndTryToPrint(code string) {
	ee, ok := lang.ptr.(PluginEvalExpression)
	if ok {
		fmt.Println(ee.EvalExpression(code))
	} else {
		lang.ptr.(PluginEval).Eval(code)
	}
}

func (lang Language) REPLLikeEval(code string) {
	rle, ok := lang.ptr.(PluginREPLLikeEval)
	if ok {
		rle.REPLLikeEval(code)
		return
	}

	ee, ok := lang.ptr.(PluginEvalExpression)
	if ok {
		fmt.Println(ee.EvalExpression(code))
	} else {
		lang.ptr.(PluginEval).Eval(code)
	}
}

func (lang Language) EvalFile(file string, args []string) {
	pef, ok := lang.ptr.(PluginEvalFile)
	if ok {
		pef.EvalFile(file, args)
	} else {
		dat, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		lang.Eval(string(dat))
	}
}

func (lang Language) REPL() {
	repl, ok := lang.ptr.(PluginREPL)
	if ok {
		repl.REPL()
	} else {
		lang.InternalREPL()
	}
}

func (lang Language) InternalREPL() {
	for {
		line, err := readline.Line(lang.ps1)
		if err != nil {
			break
		}
		lang.REPLLikeEval(line)
		readline.AddHistory(line)
	}
}

func (lang Language) SetPrompts(ps1, ps2 string) {
	ee, ok := lang.ptr.(PluginSetPrompts)
	if ok {
		ee.SetPrompts(ps1, ps2)
	}
}
