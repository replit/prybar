package main

// USING_CGO

/*
#cgo pkg-config: ruby-2.5
#include "pry_ruby.h"
*/
import "C"

import (
	"unsafe"
)

type Ruby struct {
}

func (p Ruby) Open() {
	C.pry_open()
}

func (p Ruby) Version() string {
	return C.GoString(C.pry_ruby_version())
}

func (p Ruby) Eval(code string) {
	p.EvalExpression(code)
}

func (p Ruby) EvalExpression(code string) string {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	res := C.pry_eval(ccode)
	return C.GoString(res)
}

func (p Ruby) EvalFile(file string, args []string) {
	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))
	C.pry_eval_file(cfile)
}

func (p Ruby) REPL() {
	p.Eval("require 'irb'\nbinding.irb")
}
func (p Ruby) Close() {
	C.ruby_cleanup(0)
}

var Instance = Ruby{}
