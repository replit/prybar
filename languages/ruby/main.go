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
	ps1 string
	ps2 string
}

func (p *Ruby) Open() {
	C.pry_open()
}

func (p *Ruby) Version() string {
	return C.GoString(C.pry_ruby_version())
}

func (p *Ruby) SetPrompts(ps1, ps2 string) {
	p.ps1 = ps1
	p.ps2 = ps2
}

func (p *Ruby) Eval(code string) {
	p.EvalExpression(code)
}

func (p *Ruby) EvalExpression(code string) string {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	res := C.pry_eval(ccode)
	return C.GoString(res)
}

func (p *Ruby) EvalFile(file string, args []string) {
	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))
	C.pry_eval_file(cfile)
}

func (p *Ruby) REPL() {
	cps1 := C.CString(p.ps1)
	cps1n := C.CString("$REPL_PS1")

	defer C.free(unsafe.Pointer(cps1))
	defer C.free(unsafe.Pointer(cps1n))

	cps2 := C.CString(p.ps2)
	cps2n := C.CString("$REPL_PS2")

	defer C.free(unsafe.Pointer(cps2))
	defer C.free(unsafe.Pointer(cps2n))

	vps1 := C.rb_str_new_cstr(cps1)
	vps2 := C.rb_str_new_cstr(cps2)

	C.rb_gv_set(cps1n, vps1)
	C.rb_gv_set(cps2n, vps2)

	p.Eval(`
begin
  require 'irb'
rescue LoadError
  require 'rubygems'
  gem "irb"
  require 'irb'
end

STDOUT.sync = true
IRB.setup nil
IRB.conf[:PROMPT][:PRYBAR] = {
  :PROMPT_I => $REPL_PS1,
  :PROMPT_N => $REPL_PS1,
  :PROMPT_S => $REPL_PS2,
  :PROMPT_C => $REPL_PS2,
  :RETURN => "=> %s\n"
}
IRB.conf[:PROMPT_MODE] = :PRYBAR
irb = IRB::Irb.new
irb.run IRB.conf`)
}

func (p *Ruby) Close() {
	C.ruby_cleanup(0)
}

var Instance = &Ruby{}
