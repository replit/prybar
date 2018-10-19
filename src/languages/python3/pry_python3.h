#include <Python.h>

void pry_eval_file(FILE *f, const char *file, int argn, const char *argv);
const char *pry_eval(const char *code, int start);
void pry_set_prompts(const char *ps1, const char *ps2);
void pymain_run_interactive_hook(void);
