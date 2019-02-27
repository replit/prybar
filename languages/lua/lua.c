// Copied straight from https://www.lua.org/source/5.1/lua.c.html
// with some of the cli stuff stripped out.

/*
** $Id: lua.c,v 1.160.1.2 2007/12/28 15:32:23 roberto Exp $
** Lua stand-alone interpreter
** See Copyright Notice in lua.h
*/

#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define lua_c
#define LUA_USE_LINUX

#include <lua.h>

#include <lauxlib.h>
#include <luaconf.h>
#include <lualib.h>

extern lua_State *pry_L;

static const char *progname = LUA_PROGNAME;

static void lstop(lua_State *L, lua_Debug *ar) {
  (void)ar; /* unused arg. */
  lua_sethook(L, NULL, 0, 0);
  luaL_error(L, "interrupted!");
}

static void laction(int i) {
  signal(i, SIG_DFL); /* if another SIGINT happens before lstop,
                              terminate process (default action) */
  lua_sethook(pry_L, lstop, LUA_MASKCALL | LUA_MASKRET | LUA_MASKCOUNT, 1);
}

static void l_message(const char *pname, const char *msg) {
  if (pname)
    fprintf(stderr, "%s: ", pname);
  fprintf(stderr, "%s\n", msg);
  fflush(stderr);
}

static int report(lua_State *L, int status) {
  if (status && !lua_isnil(L, -1)) {
    const char *msg = lua_tostring(L, -1);
    if (msg == NULL)
      msg = "(error object is not a string)";
    l_message(progname, msg);
    lua_pop(L, 1);
  }
  return status;
}

static int traceback(lua_State *L) {
  if (!lua_isstring(L, 1)) /* 'message' not a string? */
    return 1;              /* keep it intact */
  lua_getfield(L, LUA_GLOBALSINDEX, "debug");
  if (!lua_istable(L, -1)) {
    lua_pop(L, 1);
    return 1;
  }
  lua_getfield(L, -1, "traceback");
  if (!lua_isfunction(L, -1)) {
    lua_pop(L, 2);
    return 1;
  }
  lua_pushvalue(L, 1);   /* pass error message */
  lua_pushinteger(L, 2); /* skip this function and traceback */
  lua_call(L, 2, 1);     /* call debug.traceback */
  return 1;
}

static int docall(lua_State *L, int narg, int clear) {
  int status;
  int base = lua_gettop(L) - narg; /* function index */
  lua_pushcfunction(L, traceback); /* push traceback function */
  lua_insert(L, base);             /* put it under chunk and args */
  signal(SIGINT, laction);
  status = lua_pcall(L, narg, (clear ? 0 : LUA_MULTRET), base);
  signal(SIGINT, SIG_DFL);
  lua_remove(L, base); /* remove traceback function */
  /* force a complete garbage collection in case of errors */
  if (status != 0)
    lua_gc(L, LUA_GCCOLLECT, 0);
  return status;
}

static void print_version(void) {
  l_message(NULL, LUA_RELEASE "  " LUA_COPYRIGHT);
}

static int getargs(lua_State *L, char **argv, int n) {
  int narg;
  int i;
  int argc = 0;
  while (argv[argc])
    argc++;              /* count total number of arguments */
  narg = argc - (n + 1); /* number of arguments to the script */
  luaL_checkstack(L, narg + 3, "too many arguments to script");
  for (i = n + 1; i < argc; i++)
    lua_pushstring(L, argv[i]);
  lua_createtable(L, narg, n + 1);
  for (i = 0; i < argc; i++) {
    lua_pushstring(L, argv[i]);
    lua_rawseti(L, -2, i - n);
  }
  return narg;
}

int dostring(lua_State *L, const char *s, const char *name) {
  int status = luaL_loadbuffer(L, s, strlen(s), name) || docall(L, 0, 1);
  return report(L, status);
}

static const char *get_prompt(lua_State *L, int firstline) {
  const char *p;
  lua_getfield(L, LUA_GLOBALSINDEX, firstline ? "_PROMPT" : "_PROMPT2");
  p = lua_tostring(L, -1);
  if (p == NULL)
    p = (firstline ? LUA_PROMPT : LUA_PROMPT2);
  lua_pop(L, 1); /* remove global */
  return p;
}

static int incomplete(lua_State *L, int status) {
  if (status == LUA_ERRSYNTAX) {
    size_t lmsg;
    const char *msg = lua_tolstring(L, -1, &lmsg);
    const char *tp = msg + lmsg - (sizeof(LUA_QL("<eof>")) - 1);
    if (strstr(msg, LUA_QL("<eof>")) == tp) {
      lua_pop(L, 1);
      return 1;
    }
  }

  return 0; /* else... */
}

static int pushline(lua_State *L, int firstline) {
  char buffer[LUA_MAXINPUT];
  char *b = buffer;
  size_t l;
  const char *prmt = get_prompt(L, firstline);
  if (lua_readline(L, b, prmt) == 0)
    return 0; /* no input */
  l = strlen(b);
  if (l > 0 && b[l - 1] == '\n')            /* line ends with newline? */
    b[l - 1] = '\0';                        /* remove it */
  if (firstline && b[0] == '=')             /* first line starts with `=' ? */
    lua_pushfstring(L, "return %s", b + 1); /* change it to `return' */
  else
    lua_pushstring(L, b);
  lua_freeline(L, b);
  return 1;
}

static int try_expression(lua_State *L) {
  lua_pushliteral(L, "return ");
  lua_pushvalue(L, 1);
  lua_concat(L, 2);

  int rstatus =
      luaL_loadbuffer(L, lua_tostring(L, -1), lua_strlen(L, -1), "=stdin");

  if (rstatus) {
    lua_pop(L, 2);
    return 1;
  }

  lua_remove(L, 2);

  return 0;
}

static int loadline(lua_State *L) {
  int status = 0;
  lua_settop(L, 0);
  if (!pushline(L, 1))
    return -1; /* no input */

  for (;;) { /* repeat until gets a complete line */
    status = try_expression(L);
    if (!status)
      break;

    status = luaL_loadbuffer(L, lua_tostring(L, 1), lua_strlen(L, 1), "=stdin");

    if (!incomplete(L, status))
      break; /* cannot try to add lines? */

    if (!pushline(L, 0)) /* no more input? */
      return -1;
    lua_pushliteral(L, "\n"); /* add a new line... */
    lua_insert(L, -2);        /* ...between the two lines */
    lua_concat(L, 3);         /* join them */
  }

  lua_saveline(L, 1);
  lua_remove(L, 1); /* remove line */
  return status;
}

void dotty(lua_State *L) {
  int status;
  const char *oldprogname = progname;
  progname = NULL;
  while ((status = loadline(L)) != -1) {
    if (status == 0)
      status = docall(L, 0, 0);
    report(L, status);
    if (status == 0 && lua_gettop(L) > 0) { /* any result to print? */
      lua_getglobal(L, "print");
      lua_insert(L, 1);
      if (lua_pcall(L, lua_gettop(L) - 1, 0, 0) != 0)
        l_message(progname,
                  lua_pushfstring(L, "error calling " LUA_QL("print") " (%s)",
                                  lua_tostring(L, -1)));
    }
  }
  lua_settop(L, 0); /* clear stack */
  fputs("\n", stdout);
  fflush(stdout);
  progname = oldprogname;
}

int handle_script(lua_State *L, char **argv, int n) {
  int status;
  const char *fname;
  int narg = getargs(L, argv, n); /* collect arguments */
  lua_setglobal(L, "arg");
  fname = argv[n];
  if (strcmp(fname, "-") == 0 && strcmp(argv[n - 1], "--") != 0)
    fname = NULL; /* stdin */
  status = luaL_loadfile(L, fname);
  lua_insert(L, -(narg + 1));
  if (status == 0)
    status = docall(L, narg, 0);
  else
    lua_pop(L, narg);
  return report(L, status);
}
