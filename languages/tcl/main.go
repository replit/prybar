package main

// USING_CGO

/*
#cgo pkg-config: tcl
#include <tcl.h>
#include <stdlib.h>

*/
import "C"

import (
	"fmt"
	"os"
	"strconv"
	"unsafe"
)

type Tcl struct {
	interp *C.Tcl_Interp
}

func (p *Tcl) Open() {
	p.interp = C.Tcl_CreateInterp()

	if C.Tcl_Init(p.interp) != C.TCL_OK {
		panic("tcl interp did not init")
	}
}

func (p *Tcl) Version() string {
	major := C.int(0)
	minor := C.int(0)
	patch := C.int(0)

	C.Tcl_GetVersion(
		&major,
		&minor,
		&patch,
		nil,
	)

	return "TCL " +
		strconv.Itoa(int(major)) + "." +
		strconv.Itoa(int(minor)) + "." +
		strconv.Itoa(int(patch))
}

func (p *Tcl) eval(code string) string {
	ccode := C.CString(code)

	status := C.Tcl_Eval(p.interp, ccode)

	result := C.GoString(C.Tcl_GetStringResult(p.interp))

	if status == C.TCL_OK {
		return result
	}

	errstr := C.GoString(C.Tcl_GetStringResult(p.interp))
	fmt.Fprintf(os.Stderr, "error: %s\n", errstr)
	return ""
}

func (p *Tcl) EvalFile(file string, args []string) {
	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))

	status := C.Tcl_EvalFile(p.interp, cfile)

	if status != C.TCL_OK {
		errstr := C.GoString(C.Tcl_GetStringResult(p.interp))
		fmt.Fprintf(os.Stderr, "error: %s\n", errstr)
	}
}

func (p *Tcl) Eval(code string) {
	p.eval(code)
}

func (p *Tcl) EvalExpression(code string) string {
	return p.eval(code)
}

func (p *Tcl) Close() {
	C.Tcl_Finalize()
}
var Instance = &Tcl{}
