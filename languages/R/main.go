package main

/*
#cgo pkg-config: libr
#include <stdlib.h>
#include <Rembedded.h>
#include <RVersion.h>
#include <Rinternals.h>
#include <R_ext/Parse.h>
#include <Rembedded.h>

void pry_open() {
	const char *args[] = {"/usr/local/bin/R", "--silent", NULL};
	Rf_initEmbeddedR(2, args);
}

const char* pry_version() {
	return "R version " R_MAJOR "." R_MINOR " (" R_YEAR "-" R_MONTH "-" R_DAY ") -- \"" R_NICK "\"";
}

pry_eval(const char* code) {
	ParseStatus status;
	SEXP x = R_ParseVector(mkString(code), 1, &status, R_NilValue);
	SEXP result = eval(VECTOR_ELT(x, 0), R_GlobalEnv);
	PrintValue(result);
}

void pry_repl() {
	R_ReplDLLinit();
	while (R_ReplDLLdo1() > 0) {

    }
}

*/
import "C"

import (
	"unsafe"

	"github.com/replit/prybar/utils"
	//"strings"
)

func init() {
	utils.Register(&Language{})
}

type Language struct {
}

func (p Language) Open() {
	C.pry_open()
}

func (p Language) Version() string {
	return C.GoString(C.pry_version())
}

func (p Language) Eval(code string) {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	C.pry_eval(ccode)
}

func (p Language) REPL() {
	C.pry_repl()
}

func (p Language) Close() {
	C.Rf_endEmbeddedR(0)
}
