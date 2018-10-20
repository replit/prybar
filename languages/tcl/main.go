package main

/*
#cgo pkg-config: tcl
#include "pry_tcl.h"

*/
import "C"

import (
	"unsafe"

	"github.com/replit/prybar/utils"
)

func init() {
	utils.Register(&Language{})
}

type Language struct{}

func (p Language) Open() {
	C.pry_open()
}

func (p Language) Version() string {
	cver := C.pry_version()
	ver := C.GoString(cver)
	C.free(unsafe.Pointer(cver))
	return "Language " + ver
}

func (p Language) Eval(code string) {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	C.pry_eval(ccode)
}

func (p Language) EvalExpression(code string) string {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	return C.GoString(C.pry_eval(ccode))
}

func (p Language) Close() {
	C.pry_close()
}
