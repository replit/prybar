#include "pry_python2.h"

void pry_eval_file(FILE *f, const char *file, int argn, const char *argv)
{
    const char *xargv[argn + 1];
    const char *ptr = argv;
    for (int i = 0; i < argn; ++i)
    {
        xargv[i] = ptr;
        ptr += strlen(ptr) + 1;
    }
    xargv[argn] = NULL;
    PySys_SetArgvEx(argn, xargv, 1);
    PyRun_AnyFile(f, file);
}

const char *pry_eval(const char *code, int start)
{

    PyObject *m, *d, *s, *v;
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
    char *str = PyBytes_AS_STRING(s);
    Py_DECREF(v);
    Py_DECREF(s);
    return str;
}

void pry_set_prompts(const char *ps1, const char *ps2)
{
    PyObject *po1 = PyString_FromString(ps1);
    PyObject *po2 = PyString_FromString(ps2);
    PySys_SetObject("ps1", po1);
    PySys_SetObject("ps2", po2);
    Py_DECREF(po1);
    Py_DECREF(po2);
}
