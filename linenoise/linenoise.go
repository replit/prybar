package linenoise

/*
#include <linenoise.h>
#include <stdlib.h>

void linenoise_setup();


*/
import "C"

import (
	"fmt"
	"unsafe"
)

type Completer func(line string) []string
type Hinter func(line string) string

var completer Completer
var hinter Hinter

//export linenoiseCompletionStub
func linenoiseCompletionStub(input *C.char, completions *C.linenoiseCompletions) {
	if completer == nil {
		return
	}
	completionsSlice := completer(C.GoString(input))
	completionsLen := len(completionsSlice)
	completions.len = C.size_t(completionsLen)

	if completionsLen > 0 {
		cvec := C.malloc(C.size_t(int(unsafe.Sizeof(*(**C.char)(nil))) * completionsLen))
		cvecSlice := (*(*[999999]*C.char)(cvec))[:completionsLen]

		for i, str := range completionsSlice {
			cvecSlice[i] = C.CString(str)
		}
		completions.cvec = (**C.char)(cvec)
	}
}

//export linenoiseHintsStub
func linenoiseHintsStub(input *C.char, color *C.int, bold *C.int) *C.char {
	if hinter == nil {
		return nil
	}
	*color = 37
	*bold = 0
	return C.CString(hinter(C.GoString(input)))
}

func Linenoise(prompt string) (string, error) {
	C.linenoise_setup()
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

func LinenoiseSetCompleter(newCompleter Completer) {
	completer = newCompleter
}

func LinenoiseSetHinter(newHinter Hinter) {
	hinter = newHinter
}
