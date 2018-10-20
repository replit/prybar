#include <lua.h>

#define luai_writestring(s,l)	fwrite((s), sizeof(char), (l), stdout)
#define luai_writeline()	(luai_writestring("\n", 1), fflush(stdout))
#define LUA_PROMPT pry_prompt
#define LUA_PROMPT2 pry_prompt2

extern lua_State *pry_L;
extern char* pry_prompt;
extern char* pry_prompt2;

const char * pry_get_version();
void pry_init();
void pry_eval(const char *code);
void pry_main();
void pry_eval_file(const char* file, int argn, const char *argv);