package main

/*
#cgo pkg-config: lua-5.2
#include <stdlib.h>
#include <lua.h>
#include <lauxlib.h>
#include <lualib.h>

const char * pry_get_version() {
	return LUA_COPYRIGHT;
}

lua_State *pry_L;

void pry_init() {
	pry_L = luaL_newstate();
	luaL_openlibs(pry_L);
}

void pry_eval(const char *code) {
	luaL_dostring(pry_L, code);
}

#define luai_writestring(s,l)	fwrite((s), sizeof(char), (l), stdout)
#define luai_writeline()	(luai_writestring("\n", 1), fflush(stdout))

int lua_main (int argc, char **argv);

void pry_main() {
	char* args[2] = {"lua", NULL};
	lua_main(0, args);
}

void dotty (lua_State *L);

*/
import "C"

import (
	"unsafe"
)


type Language struct {

}

func (p Language) Open() {
	C.pry_init()
}

func (p Language) Version() string {
	return C.GoString(C.pry_get_version())
}

func (p Language) Eval(code string) {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	C.pry_eval(ccode); 
}

func (p Language) REPL() {
	fn := C.CString("<stdin>")
	defer C.free(unsafe.Pointer(fn))
	C.dotty(C.pry_L)
}

func (p Language) Close() {    
    C.lua_close(C.pry_L)
}

// exported
var Instance Language
