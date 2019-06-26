#define _GNU_SOURCE
#include <dlfcn.h>
#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/* what happens in SQLite CLI:

1. fprintf to stderr containing the "Loading resources..." message
2. printf containing version and help string

 */

// this macro lets us roll up the varargs once instead of in four separate
// places below.
#define VAR_PRINT(out, format)                                                 \
    {                                                                          \
        va_list arg;                                                           \
        va_start(arg, format);                                                 \
        int rc = internal_print(out, format, arg);                             \
        va_end(arg);                                                           \
        return rc;                                                             \
    }

/* three states:
    0: no output suppressed
    1: suppressed "Loading resources..." message on stderr
    2: suppressed version and help string on stdout
*/
int current_state = 0;

int internal_print(FILE *out, const char *format, va_list args) {
    return vfprintf(out, format, args);
}

int __fprintf_chk(FILE *out, int flag, const char *format, ...) {
    // don't do anything; pass through print
    if (current_state != 0) {
        VAR_PRINT(out, format);
    }

    // check if this is the first time we're seeing the output that we're aiming
    // to suppress
    char *expected = "-- Loading resources from /tmp/sqlite-config";
    if (getenv("PRYBAR_QUIET") && out == stderr &&
        strstr(format, expected) == 0) {
        // advance to next state and suppress output
        current_state++;
        return 0;
    }

    // this isn't the output we were looking for...
    VAR_PRINT(out, format);
}

int __printf_chk(int flag, const char *format, ...) {
    // don't do anything; pass through print
    if (current_state != 1) {
        VAR_PRINT(stdout, format);
    }

    // check if this is the first time we're seeing the output that we're aiming
    // to suppress
    char *expected = "SQLite version %s %s\nEnter\".help\" for usage "
                     "hints.\n";
    if (getenv("PRYBAR_QUIET") && strstr(format, expected) == 0) {
        // advance to next state and suppress output
        current_state++;
        return 0;
    }

    // this isn't the output we were looking for...
    VAR_PRINT(stdout, format);
}
