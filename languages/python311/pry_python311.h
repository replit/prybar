#include <Python.h>

int pry_eval_file(FILE *f, const char *file, int argn, const char *argv);
const char *pry_eval(const char *code, int start);
void pry_set_prompts(const char *ps1, const char *ps2);
int pymain_run_interactive_hook(int *exitcode);
int pymain_err_print(int *exitcode_p);
int pry_set_program_name(const char *name);
PyAPI_FUNC(int) _Py_HandleSystemExit(int *exitcode_p);
