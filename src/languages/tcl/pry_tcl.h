#include <stdlib.h>
#include <tcl.h>


void pry_open();
void pry_close();
char *pry_version();
char *pry_eval(const char *code);