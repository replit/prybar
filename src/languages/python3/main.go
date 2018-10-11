package main

/*
#cgo pkg-config: python3
#include <Python.h>

void pry_eval_file(FILE* f, const char* file, int argn, const char *argv) {
	const char* xargv[argn+1];
	const char* ptr = argv;
	for (int i = 0; i < argn; ++i) {
		xargv[i] = ptr;
		ptr += strlen(ptr) + 1;
	}
	xargv[argn] = NULL;
	PySys_SetArgvEx(argn, xargv, 1);
	PyRun_AnyFile(f, file);
}

*/
import "C"

import (
	"unsafe"
	"strings"
)


type Python struct {

}

func (p Python) Open() {
	C.Py_Initialize()
}

func (p Python) Version() string {
	return "Python " + C.GoString(C.Py_GetVersion()) + " on " + C.GoString(C.Py_GetPlatform())
}

func (p Python) Eval(code string) {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	C.PyRun_SimpleStringFlags(ccode, nil) 
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
	C.pry_eval_file(handle, cfile, C.int(len(args) + 1), argv)
}

func (p Python) REPL() {
	fn := C.CString("<stdin>")
	defer C.free(unsafe.Pointer(fn))
	C.PyRun_InteractiveLoopFlags(C.stdin, fn, nil)
}

func (p Python) Close() {    
    C.Py_Finalize()
}

// exported
var Instance Python
