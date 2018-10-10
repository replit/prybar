package main


/*
#cgo pkg-config: ruby-2.5
#include <ruby.h>
#include <ruby/version.h>

const char* prybar_ruby_version() {
	return ruby_description;
}

void prybar_eval(const char* code) {
	int state;
	VALUE result;
	result = rb_eval_string_protect(code, &state);

	if (state)
	{
		VALUE exception = rb_errinfo();
		rb_set_errinfo(Qnil);
		if (RTEST(exception)) rb_warn("%"PRIsVALUE"", rb_funcall(exception, rb_intern("full_message"), 0));
	}

}

*/
import "C"

import (
	"unsafe"
)

/*  */
type Language struct {

}

func (p Language) Open() {
	C.ruby_init()
	C.ruby_init_loadpath()
}

func (p Language) Version() string {
	//return C.GoString(C.ruby_version)
	return C.GoString(C.prybar_ruby_version())
}

func (p Language) Eval(code string) {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	C.prybar_eval(ccode)
}

func (p Language) REPL() {
	p.Eval("require 'irb'\nbinding.irb")
}

func (p Language) Close() {    
    C.ruby_cleanup(0)
}

// exported
var Instance Language
