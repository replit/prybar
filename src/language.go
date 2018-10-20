package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"
	"runtime"

	"github.com/replit/prybar/linenoise"
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

type Langauge struct {
	ptr PluginBase
}

func finalizer(f *Langauge) {
	fmt.Println("a finalizer has run.")
}

func GetLanguage(name string) *Langauge {
	base := "."
	exe, err := os.Executable()
	for {
		o, err := os.Readlink(exe)
		if err == nil {
			exe = o
		} else {
			break
		}
	}

	base = filepath.Dir(exe)
	plug, err := plugin.Open(base + "/plugins/" + name + ".so")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sym, err := plug.Lookup("Instance")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var lang PluginBase
	lang, ok := sym.(PluginBase)
	if !ok {
		fmt.Println("module did not export Language interface")
		os.Exit(1)
	}

	result := &Langauge{
		ptr: lang,
	}
	lang.Open()
	runtime.SetFinalizer(result, finalizer)
	return result
}

func (lang Langauge) Version() string {
	return lang.ptr.Version()
}

func (lang Langauge) Eval(code string) {
	lang.ptr.(PluginEval).Eval(code)
}

func (lang Langauge) EvalAndTryToPrint(code string) {
	ee, ok := lang.ptr.(PluginEvalExpression)
	if ok {
		fmt.Println(ee.EvalExpression(code))
	} else {
		lang.ptr.(PluginEval).Eval(code)
	}
}

func (lang Langauge) REPLLikeEval(code string) {
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

func (lang Langauge) EvalFile(file string, args []string) {
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

func (lang Langauge) REPL() {
	repl, ok := lang.ptr.(PluginREPL)
	if ok {
		repl.REPL()
	} else {
		lang.InternalREPL()
	}
}

func (lang Langauge) InternalREPL() {
	for {
		line, err := linenoise.Linenoise(ps1)
		if err != nil {
			break
		}
		lang.REPLLikeEval(line)
		linenoise.LinenoiseHistoryAdd(line)
	}
}

func (lang Langauge) SetPrompts(ps1, ps2 string) {
	ee, ok := lang.ptr.(PluginSetPrompts)
	if ok {
		ee.SetPrompts(ps1, ps2)
	}
}
