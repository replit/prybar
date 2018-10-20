#include "pry_tcl.h"

Tcl_Interp *interp;

void pry_open()
{
    interp = Tcl_CreateInterp();
    if (Tcl_Init(interp) != TCL_OK)
    {
        exit(8);
    }
}

void pry_close()
{
    Tcl_Finalize();
}

char *pry_version()
{
    char *result = malloc(100);
    int a, b, c;
    Tcl_GetVersion(&a, &b, &c, NULL);
    sprintf(result, "%d.%d.%d", a, b, c);
    return result;
}

char *pry_eval(const char *code)
{
    if (Tcl_Eval(interp, code) == TCL_OK)
    {
        return Tcl_GetStringResult(interp);
    }
    else
    {
        fprintf(stderr, "error: %s\n", Tcl_GetStringResult(interp));
        return "";
    }
}