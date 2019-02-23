#include "pry_ruby.h"

VALUE binding;

void pry_open() {
  ruby_init();
  ruby_init_loadpath();
  binding = rb_binding_new();
}

const char *pry_ruby_version() { return ruby_description; }

char *pry_eval(const char *code) {
  int state = 0;

  VALUE result = rb_eval_string_wrap(code, &state);

  if (state) {
    VALUE exception = rb_errinfo();
    rb_set_errinfo(Qnil);
    if (RTEST(exception))
      rb_warn(
          "%" PRIsVALUE "",
          rb_funcall(exception, rb_intern("full_message"), 0));
    return NULL;
  }

  VALUE str = rb_sprintf("%" PRIsVALUE "", result);
  return StringValueCStr(str);
}

void pry_eval_file(char *file) {
  char *options[] = {"ruby", file};
  void *node = ruby_options(2, options);

  int state = 0;
  if (ruby_executable_node(node, &state)) {
    ruby_exec_node(node);
  }
}
