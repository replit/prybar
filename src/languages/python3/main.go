package main

/*
#cgo pkg-config: python3
#include <Python.h>

void pry_eval_file(FILE* f, const char* file, int argn, const char *argv) {
	wchar_t* xargv[argn+1];
	const char* ptr = argv;
	for (int i = 0; i < argn; ++i) {
		xargv[i] = Py_DecodeLocale(ptr, NULL);
		ptr += strlen(ptr) + 1;
	}
	xargv[argn] = NULL;
	PySys_SetArgvEx(argn, xargv, 1);
	PyRun_AnyFile(f, file);
}

const char* pry_eval(const char *code, int start) {

	PyObject *m, *d, *s, *v;
	PyObject *c;
	m = PyImport_AddModule("__main__");

	if (m == NULL) return NULL;

	d = PyModule_GetDict(m);
	c = Py_CompileString(code, "(eval)", start);
	if (c == NULL) {
		PyErr_Print();
		return NULL;
	}
	v = PyEval_EvalCode(c, d, d);
	if (v == NULL) {
		PyErr_Print();
		return NULL;
	}
	s = PyObject_Str(v);
	if (s == NULL) {
		PyErr_Print();
		return NULL;
	}
	char *str = PyBytes_AS_STRING(s);
	Py_DECREF(v);
	Py_DECREF(s);
	return str;
}

void pry_set_prompts(const char* ps1, const char* ps2) {
	PyObject *po1 = PyUnicode_FromString(ps1);
	PyObject *po2 = PyUnicode_FromString(ps2);
	PySys_SetObject("ps1", po1);
	PySys_SetObject("ps2", po2);
	Py_DECREF(po1);
	Py_DECREF(po2);
}

//From python3 sourcecode
static void
pymain_run_interactive_hook(void)
{
    PyObject *sys, *hook, *result;
    sys = PyImport_ImportModule("sys");
    if (sys == NULL) {
        goto error;
    }

    hook = PyObject_GetAttrString(sys, "__interactivehook__");
    Py_DECREF(sys);
    if (hook == NULL) {
        PyErr_Clear();
        return;
    }

    result = _PyObject_CallNoArg(hook);
    Py_DECREF(hook);
    if (result == NULL) {
        goto error;
    }
    Py_DECREF(result);

    return;

error:
    PySys_WriteStderr("Failed calling sys.__interactivehook__\n");
    PyErr_Print();
}

*/
import "C"

import (
	"strings"
	"unsafe"
)

type Python struct {
}

func (p Python) Open() {
	C.Py_Initialize()
}

func (p Python) Version() string {
	return "Python " + C.GoString(C.Py_GetVersion()) + " on " + C.GoString(C.Py_GetPlatform())
}

func (p Python) Eval(code string) {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	C.GoString(C.pry_eval(ccode, C.Py_file_input))
}

func (p Python) EvalExpression(code string) string {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	return C.GoString(C.pry_eval(ccode, C.Py_eval_input))
}

func (p Python) EvalFile(file string, args []string) {
	handle := C.stdin
	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))

	if file != "-" {
		cmode := C.CString("r")

		defer C.free(unsafe.Pointer(cmode))
		handle = C.fopen(cfile, cmode)
		defer C.fclose(handle)
	}

	argv := C.CString(file + "\x00" + strings.Join(args, "\x00"))
	defer C.free(unsafe.Pointer(argv))
	C.pry_eval_file(handle, cfile, C.int(len(args)+1), argv)
}

func (p Python) REPLLikeEval(code string) {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	C.pry_eval(ccode, C.Py_single_input)
}

func (p Python) REPL() {
	C.pymain_run_interactive_hook()

	fn := C.CString("<stdin>")
	defer C.free(unsafe.Pointer(fn))
	C.PyRun_InteractiveLoopFlags(C.stdin, fn, nil)
}

func (p Python) SetPrompts(ps1, ps2 string) {
	cps1 := C.CString(ps1)
	defer C.free(unsafe.Pointer(cps1))
	cps2 := C.CString(ps2)
	defer C.free(unsafe.Pointer(cps2))

	C.pry_set_prompts(cps1, cps2)
}

func (p Python) Close() {
	C.Py_Finalize()
}

// exported
var Instance Python
