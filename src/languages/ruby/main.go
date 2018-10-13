package main


/*
#cgo pkg-config: ruby-2.5
#include <ruby.h>
#include <ruby/version.h>

VALUE binding;

void pry_open() {
	ruby_init();
	ruby_init_loadpath();
	binding = rb_binding_new();
}

const char* pry_ruby_version() {
	return ruby_description;
}



char * pry_eval(const char* code) {
	
	int state;
	VALUE result;
	result = rb_eval_string_wrap(code, &state);

	if (state)
	{
		VALUE exception = rb_errinfo();
		rb_set_errinfo(Qnil);
		if (RTEST(exception)) rb_warn("%"PRIsVALUE"", rb_funcall(exception, rb_intern("full_message"), 0));
		return NULL;
	} else {
		VALUE str = rb_sprintf("%"PRIsVALUE"", result);
		return StringValueCStr(str);
	}

}

void pry_eval_file(const char* file ) {
	char* options[] = { "ruby", file };
	void* node = ruby_options(2, options);

	int state;
	if (ruby_executable_node(node, &state))
	{
		state = ruby_exec_node(node);
	}
}

*/
import "C"

import (
	"unsafe"
)

/*  */
type Ruby struct {

}

func (p Ruby) Open() {
	C.pry_open();
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
	return C.GoString(res);
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

// exported
var Instance Ruby
