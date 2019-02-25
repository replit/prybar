package main

// USING_CGO

/*
#cgo pkg-config: tcl
#include "pry_tcl.h"

*/
import "C"

import (
	"unsafe"
)

type Tcl struct{}

func (p Tcl) Open() {
	C.pry_open()
}

func (p Tcl) Version() string {
	cver := C.pry_version()
	ver := C.GoString(cver)
	C.free(unsafe.Pointer(cver))
	return "TCL " + ver
}

func (p Tcl) Eval(code string) {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	C.pry_eval(ccode)
}

func (p Tcl) EvalExpression(code string) string {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	return C.GoString(C.pry_eval(ccode))
}

func (p Tcl) Close() {
	C.pry_close()
}

var Instance = Tcl{}
