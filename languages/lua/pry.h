#include <lua.h>

extern lua_State *pry_L;

const char *pry_get_version(void);
void pry_eval(const char *code);
void pry_eval_file(char *file);
void pry_do_repl(void);

// from the lua repl lib (lua.c)
void dotty(lua_State *L);
int handle_script(lua_State *L, char **argv, int n);
int dostring(lua_State *L, const char *s, const char *name);
