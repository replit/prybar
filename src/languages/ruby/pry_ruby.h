#include <ruby.h>
#include <ruby/version.h>

void pry_open();
const char *pry_ruby_version();
char *pry_eval(const char *code);
void pry_eval_file(const char *file);

