#include <string.h>

#include <lauxlib.h>
#include <lualib.h>

#include "pry.h"

char* pry_prompt = "> ";
char* pry_prompt2 = ">> ";
lua_State *pry_L;

const char * pry_get_version() {
	return LUA_COPYRIGHT;
}


void pry_init() {
	pry_L = luaL_newstate();
	luaL_openlibs(pry_L);
}

void pry_eval(const char *code) {
	luaL_dostring(pry_L, code);
}


int lua_main (int argc, char **argv);

void pry_main() {
	char* args[2] = {"lua", NULL};
	lua_main(0, args);
}

int handle_script (lua_State *L, char **argv, int n);
void pry_eval_file(const char* file, int argn, const char *argv) {
	int status, result;
	const char* xargv[argn+1];
	const char* ptr = argv;
	for (int i = 0; i < argn; ++i) {
		xargv[i] = ptr;
		ptr += strlen(ptr) + 1;
	}
	xargv[argn] = NULL;

	handle_script(pry_L, xargv, 0);
}