package main

// USING_CGO

/*
#cgo pkg-config: lua-5.1
#cgo LDFLAGS: -lreadline
#include <stdlib.h>
#include <lua.h>
#include <lauxlib.h>
#include <lualib.h>
#include "pry.h"

*/
import "C"

import (
	"unsafe"
)

type Lua struct{}

func (p Lua) Open() {
	C.pry_L = C.luaL_newstate()
	C.luaL_openlibs(C.pry_L)
}

func (p Lua) Version() string {
	return C.GoString(C.pry_get_version())
}

func (p Lua) SetPrompts(ps1, ps2 string) {
	c1 := C.CString(ps1)
	C.lua_pushstring(C.pry_L, c1)
	C.free(unsafe.Pointer(c1))

	C.lua_setfield(C.pry_L, C.LUA_GLOBALSINDEX, C.CString("_PROMPT"))

	c2 := C.CString(ps2)
	C.lua_pushstring(C.pry_L, c2)
	C.free(unsafe.Pointer(c2))

	C.lua_setfield(C.pry_L, C.LUA_GLOBALSINDEX, C.CString("_PROMPT2"))
}

func (p Lua) Eval(code string) {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	C.pry_eval(ccode)
}

func (p Lua) EvalFile(file string) {
	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))

	C.pry_eval_file(cfile)
}

func (p Lua) REPL() {
	C.pry_do_repl()
}

func (p Lua) Close() {
	C.lua_close(C.pry_L)
}

var Instance = Lua{}
