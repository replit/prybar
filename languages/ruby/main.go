package main

/*
#cgo pkg-config: ruby-2.5
#include "pry_ruby.h"
*/
import "C"

import (
	"unsafe"

	"github.com/replit/prybar/utils"
)

func init() {
	utils.Register(&Language{})
}

/*  */
type Language struct {
}

func (p Language) Open() {
	C.pry_open()
}

func (p Language) Version() string {
	return C.GoString(C.pry_ruby_version())
}

func (p Language) Eval(code string) {
	p.EvalExpression(code)
}

func (p Language) EvalExpression(code string) string {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	res := C.pry_eval(ccode)
	return C.GoString(res)
}

func (p Language) EvalFile(file string, args []string) {
	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))
	C.pry_eval_file(cfile)
}

func (p Language) REPL() {
	p.Eval("require 'irb'\nbinding.irb")
}
func (p Language) Close() {
	C.ruby_cleanup(0)
}
