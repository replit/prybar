package main

// USING_CGO

/*
#cgo pkg-config: python-3.11-embed
#include "pry_python311.h"
*/
import "C"

import (
	"fmt"
	"os"
	"path"
	"strings"
	"unicode/utf16"
	"unsafe"
)

var programName *C.wchar_t

type Python struct{}

func (p Python) Open() {
	name := "python"
	virtualEnv, virtualEnvSet := os.LookupEnv("VIRTUAL_ENV")
	if virtualEnvSet {
		name = path.Join(virtualEnv, "bin", "python")
	}

	wname := utf16.Encode([]rune(name))

	var pyConfig C.PyConfig
	C.PyConfig_InitPythonConfig(&pyConfig)
	res := C.PyConfig_SetString(&pyConfig, &pyConfig.program_name, (*C.Wchar_t)(wname))
	if C.PyStatus_Exception(res) != 0 {
		C.PyConfig_Clear(&pyConfig)
		panic(fmt.Sprintf("cannot set prybar-python311 program name to '%s'", name))
	}

	C.Py_InitializeFromConfig(&pyConfig)
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
	handle := C.stdin
	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))

	if file != "-" {
		cmode := C.CString("r")

		defer C.free(unsafe.Pointer(cmode))
		handle = C.fopen(cfile, cmode)
		defer C.fclose(handle)
	}

	argv := C.CString(file + "\x00" + strings.Join(args, "\x00"))
	defer C.free(unsafe.Pointer(argv))

	status := (C.pry_eval_file(handle, cfile, C.int(len(args)+1), argv))

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
