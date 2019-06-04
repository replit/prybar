#include "pry_python3.h"

void pry_eval_file(FILE *f, const char *file, int argn, const char *argv)
{
    wchar_t *xargv[argn + 1];
    const char *ptr = argv;
    for (int i = 0; i < argn; ++i)
    {
        xargv[i] = Py_DecodeLocale(ptr, NULL);
        ptr += strlen(ptr) + 1;
    }
    xargv[argn] = NULL;
    PySys_SetArgvEx(argn, xargv, 1);
    PyRun_AnyFile(f, file);

}

const char *pry_eval(const char *code, int start)
{

    PyObject *m, *d, *s, *t, *v;
    PyObject *c;
    m = PyImport_AddModule("__main__");

    if (m == NULL)
        return NULL;

    d = PyModule_GetDict(m);
    c = Py_CompileString(code, "(eval)", start);
    if (c == NULL)
    {
        PyErr_Print();
        return NULL;
    }
    v = PyEval_EvalCode(c, d, d);
    if (v == NULL)
    {
        PyErr_Print();
        return NULL;
    }
    s = PyObject_Str(v);
    if (s == NULL)
    {
        PyErr_Print();
        return NULL;
    }

    t = PyUnicode_AsUTF8String(s);
    if (t == NULL)
    {
        PyErr_Print();
        return NULL;
    }


    char *str = PyBytes_AsString(t);
    Py_DECREF(v);
    Py_DECREF(s);
    Py_DECREF(t);
    return str;
}

void pry_set_prompts(const char *ps1, const char *ps2)
{
    PyObject *po1 = PyUnicode_FromString(ps1);
    PyObject *po2 = PyUnicode_FromString(ps2);
    PySys_SetObject("ps1", po1);
    PySys_SetObject("ps2", po2);
    Py_DECREF(po1);
    Py_DECREF(po2);
}

//From python3 sourcecode
void pymain_run_interactive_hook(void)
{
    PyObject *sys, *hook, *result;
    sys = PyImport_ImportModule("sys");
    if (sys == NULL)
    {
        goto error;
    }

    hook = PyObject_GetAttrString(sys, "__interactivehook__");
    Py_DECREF(sys);
    if (hook == NULL)
    {
        PyErr_Clear();
        return;
    }

    result = _PyObject_CallNoArg(hook);
    Py_DECREF(hook);
    if (result == NULL)
    {
        goto error;
    }
    Py_DECREF(result);

    return;

error:
    PySys_WriteStderr("Failed calling sys.__interactivehook__\n");
    PyErr_Print();
}
