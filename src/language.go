package main

import(
	"os"
	"plugin"
	"fmt"
	"runtime"
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

type PluginEvalFile interface {
	PluginBase
	EvalFile(file string, args []string)
}

type PluginREPL interface {
	PluginBase
	REPL()
}

type Langauge struct {
	ptr PluginBase
}

func finalizer(f *Langauge) {
        fmt.Println("a finalizer has run.")
} 

func GetLanguage(name string) *Langauge {	
	plug, err := plugin.Open("./plugins/" + name + ".so")

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

	result := &Langauge {
		ptr: lang,
	}
	lang.Open();
	runtime.SetFinalizer(result, finalizer)
	return result
}

func (lang Langauge) Version() string {
	return lang.ptr.Version()
}

func (lang Langauge) Eval(code string) {
	lang.ptr.(PluginEval).Eval(code)
}

func (lang Langauge) EvalFile(file string, args []string) {
	lang.ptr.(PluginEvalFile).EvalFile(file, args)
}

func (lang Langauge) REPL() {
	lang.ptr.(PluginREPL).REPL()
}
