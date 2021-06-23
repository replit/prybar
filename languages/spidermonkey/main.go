package main

// USING_CGO

/*
#cgo pkg-config: mozjs185
#include "jsapi.h"

static JSClass global_class = {
	"global", JSCLASS_GLOBAL_FLAGS,
	JS_PropertyStub, JS_PropertyStub, JS_PropertyStub, JS_StrictPropertyStub,
	JS_EnumerateStub, JS_ResolveStub, JS_ConvertStub, JS_FinalizeStub,
	JSCLASS_NO_OPTIONAL_MEMBERS
};

JSRuntime *rt;
JSContext *cx;
JSObject  *global;

int pry_open() {


	rt = JS_NewRuntime(8 * 1024 * 1024);
	if (rt == NULL) return 1;

	cx = JS_NewContext(rt, 8192);
	if (cx == NULL) return 1;

	JS_SetOptions(cx, JSOPTION_VAROBJFIX | JSOPTION_JIT | JSOPTION_METHODJIT);
	JS_SetVersion(cx, JSVERSION_LATEST);
	//JS_SetErrorReporter(cx, reportError);


	global = JS_NewCompartmentAndGlobalObject(cx, &global_class, NULL);
	if (global == NULL) return 1;


    if (!JS_InitStandardClasses(cx, global)) return 1;
    return 0;
}

int pry_eval(const char* script) {
    jsval rval;
    JSString *str;
    JSBool ok;
    const char *filename = "noname";
    uintN lineno = 0;

    ok = JS_EvaluateScript(cx, global, script, strlen(script), filename, lineno, &rval);
    if (rval == JS_FALSE)
    return 1;

    str = JS_ValueToString(cx, rval);
    printf("%s\n", JS_EncodeString(cx, str));
}

void pry_close() {
	JS_DestroyContext(cx);
	JS_DestroyRuntime(rt);
	JS_ShutDown();
}

*/
import "C"

import (
	"unsafe"
)

type SpiderMonkey struct{}

func (p SpiderMonkey) Open() {
	C.pry_open()
}

func (p SpiderMonkey) Version() string {
	return C.GoString(C.JS_GetImplementationVersion())
}

func (p SpiderMonkey) Eval(code string) {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	C.pry_eval(ccode)
}

func (p SpiderMonkey) Close() {
	C.pry_close()
}

var Instance = SpiderMonkey{}
