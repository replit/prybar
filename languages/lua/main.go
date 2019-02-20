package main

// USING_CGO

/*
#cgo pkg-config: lua-5.2
#include <stdlib.h>
#include <lua.h>
#include "pry.h"

int lua_main (int argc, char **argv);
void dotty (lua_State *L);
int pmain (lua_State *L);



*/
import "C"

import (
	"strings"
	"unsafe"
)

type Lua struct{}

func (p Lua) Open() {
	C.pry_init()
}

func (p Lua) Version() string {
	return C.GoString(C.pry_get_version())
}

func (p Lua) Eval(code string) {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	C.pry_eval(ccode)
}

func (p Lua) EvalFile(file string, args []string) {
	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))

	argv := C.CString(file + "\x00" + strings.Join(args, "\x00"))
	defer C.free(unsafe.Pointer(argv))
	C.pry_eval_file(cfile, C.int(len(args)+1), argv)
}

func (p Lua) REPL() {
	fn := C.CString("<stdin>")
	defer C.free(unsafe.Pointer(fn))
	C.dotty(C.pry_L)
}

func (p Lua) Close() {
	C.lua_close(C.pry_L)
}

var Instance = Lua{}
