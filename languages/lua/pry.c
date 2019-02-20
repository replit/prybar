#include <string.h>

#include <lauxlib.h>
#include <lualib.h>

#include "pry.h"

lua_State *pry_L;

const char *pry_get_version(void) { return LUA_RELEASE "  " LUA_COPYRIGHT; }

void pry_eval(const char *code) { dostring(pry_L, code, "<eval>"); }

void pry_do_repl(void) { dotty(pry_L); }

void pry_eval_file(char *file) { handle_script(pry_L, &file, 0); }
