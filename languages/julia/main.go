package main

// USING_CGO

/*
#cgo CFLAGS: -I/usr/include/julia
#cgo LDFLAGS: -ljulia
#include "pry.h"
*/
import "C"

import (
	"unsafe"
)

type Julia struct{}

func (p Julia) Open() {
	C.setup()
}

func (p Julia) SetPrompts(ps1, ps2 string) {
	cps1 := C.CString(ps1)
	defer C.free(unsafe.Pointer(cps1))

	C.set_prompt(cps1)
}

func (p Julia) Version() string {
	return C.GoString(C.get_banner())
}

func (p Julia) Eval(code string) {
	cstr := C.CString(code)
	defer C.free(unsafe.Pointer(cstr))

	C.eval(cstr)
}

func (p Julia) EvalFile(file string, args []string) {
	cstr := C.CString(file)
	C.eval_file(cstr)
	C.free(unsafe.Pointer(cstr))
}

func (p Julia) REPL() {
	C.run_repl()
}

func (p Julia) Close() {
	C.cleanup()
}

var Instance = Julia{}
