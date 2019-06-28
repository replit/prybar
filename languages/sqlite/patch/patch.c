#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/* a shared library designed to be LD_PRELOADed into SQLite CLI

If PRYBAR_QUIET is set in the environment, then this suppresses the "Loading
resources" message due to the CLI's `-init` flag, and also hides the version
header and help message.

What happens in SQLite CLI:

1. fprintf to stderr containing the "Loading resources..." message
2. printf containing version and help string

We need to suppress the output if the PRYBAR_QUIET env var is set. So, we wait
for the CLI to make calls to libc print functions and intercept them. If they're
what we expect, then we don't print anything and then advance the state as
necessary.

Check out the SQLite CLI code to see what it does:
https://www.sqlite.org/src/artifact/4e1bcf8c70b8fb97

*/

// this macro lets us roll up the varargs once instead of in four separate
// places below.
#define VAR_PRINT(out, format)                                                 \
    {                                                                          \
        va_list arg;                                                           \
        va_start(arg, format);                                                 \
        int rc = vfprintf(out, format, arg);                                   \
        va_end(arg);                                                           \
        return rc;                                                             \
    }

/* three states:
    - no output suppressed
    - suppressed "Loading resources..." message on stderr
    - suppressed version and help string on stdout
*/
static enum {
    INITIAL,
    LOADING_SUPPRESSED,
    VERSION_SUPPRESSED
} current_state = INITIAL;

int __fprintf_chk(FILE *out, int flag, const char *format, ...) {
    // don't do anything; pass through print
    if (current_state != INITIAL) {
        VAR_PRINT(out, format);
    }

    // check if this is the first time we're seeing the output that we're aiming
    // to suppress
    const char *expected = "-- Loading resources from %s\n";
    if (getenv("PRYBAR_QUIET") && out == stderr &&
        strcmp(format, expected) == 0) {
        // advance to next state and suppress output
        current_state = LOADING_SUPPRESSED;
        return 0;
    }

    // this isn't the output we were looking for...
    VAR_PRINT(out, format);
}

int __printf_chk(int flag, const char *format, ...) {
    // don't do anything; pass through print
    if (current_state != LOADING_SUPPRESSED) {
        VAR_PRINT(stdout, format);
    }

    // check if this is the first time we're seeing the output that we're aiming
    // to suppress
    const char *expected = "SQLite version %s %.19s\nEnter \".help\" for usage "
                           "hints.\n";
    if (getenv("PRYBAR_QUIET") && strcmp(format, expected) == 0) {
        // advance to next state and suppress output
        current_state = VERSION_SUPPRESSED;
        return 0;
    }

    // this isn't the output we were looking for...
    VAR_PRINT(stdout, format);
}
