package main

/*
#cgo pkg-config: tcl
#include <stdlib.h>
#include <tcl.h>

Tcl_Interp *interp;
int ExtendTcl (Tcl_Interp *interp) {

    return TCL_OK;
}

void pry_open() {
	interp = Tcl_CreateInterp();
	if (Tcl_Init(interp) != TCL_OK) {
        exit(8);
    }
}

void pry_close() {
	Tcl_Finalize();
}

char * pry_version() {
	char *result = malloc(100);
	int a, b, c;
	Tcl_GetVersion(&a, &b, &c, NULL);
	sprintf(result, "%d.%d.%d",a,b,c);
	return result;
}

char * pry_eval(const char* code) {
	if ( Tcl_Eval(interp, code) == TCL_OK ) {
		return Tcl_GetStringResult(interp);
	} else {
		fprintf(stderr, "error: %s\n", Tcl_GetStringResult (interp));
		return "";
	}
}

*/
import "C"

import (
	"unsafe"
)


type Python struct {

}

func (p Python) Open() {
	C.pry_open()
}

func (p Python) Version() string {
	cver := C.pry_version()
	ver := C.GoString(cver)
	C.free(unsafe.Pointer(cver))
	return "TCL " + ver
}

func (p Python) Eval(code string) {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	C.pry_eval(ccode)
}

func (p Python) EvalExpression(code string) string {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	return C.GoString(C.pry_eval(ccode))
}


func (p Python) Close() {    
	C.pry_close()
}

// exported
var Instance Python
