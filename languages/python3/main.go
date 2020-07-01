package main

// USING_CGO

/*
#cgo pkg-config: python-3.8-embed
#include "pry_python3.h"
*/
import "C"

import (
	"fmt"
	"os"
	"path"
	"strings"
	"unsafe"
)

var programName *C.wchar_t

func Py_SetProgramName(name string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	newProgramName := C.Py_DecodeLocale(cname, nil)
	if newProgramName == nil {
		return fmt.Errorf("fail to call Py_DecodeLocale on '%s'", name)
	}
	C.Py_SetProgramName(newProgramName)

	//no operation is performed if nil
	C.PyMem_RawFree(unsafe.Pointer(programName))
	programName = newProgramName

	return nil
}

func init() {
	name := "python"
	virtualEnv, virtualEnvSet := os.LookupEnv("VIRTUAL_ENV")
	if virtualEnvSet {
		name = path.Join(virtualEnv, "bin", "python")
	}
	err := Py_SetProgramName(name)
	if err != nil {
		panic(fmt.Sprintf("cannot set prybar-python3 program name to '%s': %s", name, err))
	}
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

func (p Python) EvalFile(file string, args []string) {
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
	C.pry_eval_file(handle, cfile, C.int(len(args)+1), argv)
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
	C.pymain_run_interactive_hook()

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
