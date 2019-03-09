// This file is a part of Julia. License is MIT: https://julialang.org/license

/*
  repl.c
  system startup, main(), and console interaction
*/

#define JULIA_ENABLE_THREADING

#include <assert.h>
#include <ctype.h>
#include <errno.h>
#include <inttypes.h>
#include <limits.h>
#include <math.h>
#include <setjmp.h>
#include <signal.h>
#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <time.h>

#include <julia.h>

JULIA_DEFINE_FAST_TLS()

void fancy_repl(void) {
  assert(jl_base_module && "oh no, we have no repl?");

  jl_function_t *start_client =
      (jl_function_t *)jl_get_global(jl_base_module, jl_symbol("_start"));

  JL_TRY {
    size_t last_age = jl_get_ptls_states()->world_age;
    jl_get_ptls_states()->world_age = jl_get_world_counter();

    jl_apply(&start_client, 1);

    jl_get_ptls_states()->world_age = last_age;
  }
  JL_CATCH { jl_no_exc_handler(jl_current_exception()); }
}

static int exec_program(const char *program) {
  JL_TRY { jl_load(jl_main_module, program); }
  JL_CATCH {
    jl_value_t *errs = jl_stderr_obj();
    volatile int shown_err = 0;
    jl_printf(JL_STDERR, "error during bootstrap:\n");
    JL_TRY {
      if (errs) {
        jl_value_t *showf = jl_get_function(jl_base_module, "show");
        if (showf != NULL) {
          jl_call2(showf, errs, jl_current_exception());
          jl_printf(JL_STDERR, "\n");
          shown_err = 1;
        }
      }
    }
    JL_CATCH {}
    if (!shown_err) {
      jl_static_show(JL_STDERR, jl_current_exception());
      jl_printf(JL_STDERR, "\n");
    }
    jlbacktrace();
    jl_printf(JL_STDERR, "\n");
    return 1;
  }
  return 0;
}

const char *get_banner(void) {
  jl_value_t *ret = jl_eval_string("io = IOBuffer(); Base.banner(IOContext(io, "
                                   ":color => true)); String(take!(io))");

  assert(jl_typeis(ret, jl_string_type) && "banner should be a string");

  return jl_string_ptr(ret);
}

void setup() {
  jl_init();
  libsupport_init();

  jl_options.banner = 0;
  jl_options.color = 1;

  jl_eval_string("import REPL");
  jl_set_global(jl_base_module, jl_symbol("replit_prompt"),
                jl_cstr_to_string("julia> "));
  jl_eval_string(
      "atreplinit(function (r)"
      "pparts = match(r\"(\\e\\[.+m)?(.*)\", Base.replit_prompt);"
      "r.prompt_color = pparts[1] == nothing ? \"\" : String(pparts[1]);"
      "r.interface = REPL.setup_interface(r, true, r.options.extra_keymap); "
      "r.interface.modes[1].prompt ="
      " pparts[2] == nothing ? \"\" : String(pparts[2]);"
      "r.interface.modes[2].prompt = \"shell> \";"
      "r.interface.modes[3].prompt = \"help?> \";"
      "end)");
}

void cleanup() { jl_atexit_hook(0); }

void eval(const char *str) { jl_eval_string(str); }

void eval_file(const char *path) { exec_program(path); }

void run_repl() { fancy_repl(); }

void set_prompt(const char *prompt) {
  jl_set_global(jl_base_module, jl_symbol("replit_prompt"),
                jl_cstr_to_string(prompt));
}
