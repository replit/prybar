package main

/*
#cgo pkg-config: tcl
#include "pry_tcl.h"

*/
import "C"

import (
	"unsafe"
)

type Python struct {
}

func (p Python) Open() {
	C.pry_open()
}

func (p Python) Version() string {
	cver := C.pry_version()
	ver := C.GoString(cver)
	C.free(unsafe.Pointer(cver))
	return "TCL " + ver
}

func (p Python) Eval(code string) {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	C.pry_eval(ccode)
}

func (p Python) EvalExpression(code string) string {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	return C.GoString(C.pry_eval(ccode))
}

func (p Python) Close() {
	C.pry_close()
}

// exported
var Instance Python
