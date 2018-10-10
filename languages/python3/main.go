package main

/*
#cgo pkg-config: python3
#include <Python.h>
*/
import "C"

import (
	"unsafe"
)


type Language struct {

}

func (p Language) Open() {
	C.Py_Initialize()
}

func (p Language) Version() string {
	return C.GoString(C.Py_GetVersion())
}

func (p Language) Eval(code string) {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	C.PyRun_SimpleStringFlags(ccode, nil) 
}

func (p Language) REPL() {
	fn := C.CString("<stdin>")
	defer C.free(unsafe.Pointer(fn))
	C.PyRun_InteractiveLoopFlags(C.stdin, fn, nil)
}

func (p Language) Close() {    
    C.Py_Finalize()
}

// exported
var Instance Language
