package main

// USING_CGO

/*
#cgo pkg-config: python-3.11-embed
#include "pry_python311.h"
*/
import "C"

import (
	"os"
	"path"
	"strings"
	"unsafe"
)

func GetProgramName() string {
	name := "python"
	virtualEnv, virtualEnvSet := os.LookupEnv("VIRTUAL_ENV")
	if virtualEnvSet {
		name = path.Join(virtualEnv, "bin", "python")
	}
	return name
}

type Python struct{}

func (p Python) Open() {
	C.Py_Initialize()
	p.loadModule("readline")
	p.Eval("import signal")
	p.Eval("signal.signal(signal.SIGINT, signal.default_int_handler)")
	p.Eval("del signal")
}

func (p Python) Version() string {
	return "Python " + C.GoString(C.Py_GetVersion()) + " on " + C.GoString(C.Py_GetPlatform())
}

func (p Python) Eval(code string) {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	C.GoString(C.pry_eval(ccode, C.Py_file_input))
}

func (p Python) EvalExpression(code string) string {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	return C.GoString(C.pry_eval(ccode, C.Py_eval_input))
}

func (p Python) EvalFile(file string, args []string) int {
	argv := C.CString(file + "\x00" + strings.Join(args, "\x00"))
	defer C.free(unsafe.Pointer(argv))

	programName := GetProgramName()
	cprogramName := C.CString(programName)
	defer C.free(unsafe.Pointer(cprogramName))

	status := (C.pry_eval_file(cprogramName, C.int(len(args)+1), argv))

	// if status is non-zero an error occured.
	if status != 0 {
		return 1
	}

	return 0
}

func (p Python) REPLLikeEval(code string) {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	C.pry_eval(ccode, C.Py_single_input)
}

func (p Python) loadModule(mod string) {
	cmode := C.CString(mod)
	defer C.free(unsafe.Pointer(cmode))
	C.PyImport_ImportModule(cmode)
}

func (p Python) REPL() {
	var exitCode C.int
	if C.pymain_run_interactive_hook(&exitCode) != 0 {
		os.Exit(int(exitCode))
	}

	fn := C.CString("<stdin>")
	defer C.free(unsafe.Pointer(fn))
	C.PyRun_InteractiveLoopFlags(C.stdin, fn, nil)
}

func (p Python) SetPrompts(ps1, ps2 string) {
	cps1 := C.CString(ps1)
	defer C.free(unsafe.Pointer(cps1))
	cps2 := C.CString(ps2)
	defer C.free(unsafe.Pointer(cps2))

	C.pry_set_prompts(cps1, cps2)
}

func (p Python) Close() {
	C.Py_Finalize()
}

var Instance = Python{}
