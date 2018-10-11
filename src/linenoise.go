package main

/*
#include <stdlib.h>
#include <linenoise.h>
*/
import "C"

import (
	"unsafe"
	"fmt"
)

func Linenoise(prompt string) (string,error) {
	cprompt := C.CString(prompt)
	defer C.free(unsafe.Pointer(cprompt))
	cresult := C.linenoise(cprompt)
	if cresult == nil {
		return "", fmt.Errorf("No input")
	}
	result := C.GoString(cresult)
	C.linenoiseFree(unsafe.Pointer(cresult))
	return result, nil
}

func LinenoiseHistoryAdd(line string) {
	data := C.CString(line)
	defer C.free(unsafe.Pointer(data))
	C.linenoiseHistoryAdd(data)
}