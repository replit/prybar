#define _GNU_SOURCE
#include <dlfcn.h>
#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int (*orig_fprintf_chk)(FILE *out, int flag, const char *format, ...) = NULL;

int quiet;

int __printf_chk(int flag, const char *format, ...) {
    va_list arg;
    va_start(arg, format);
    int rc = __fprintf_chk(stdout, flag, format);
    va_end(arg);
    return rc;
}

int __fprintf_chk(FILE *out, int flag, const char *format, ...) {
    if (!orig_fprintf_chk) {
        orig_fprintf_chk = dlsym(RTLD_NEXT, "__fprintf_chk");

        if (getenv("PRYBAR_QUIET")) {
            quiet = 1;
        } else {
            quiet = 0;
        }
    }

    if (quiet) {
        return 0;
    }

    va_list arg;
    va_start(arg, format);
    int rc = fprintf(out, format, arg);
    va_end(arg);

    // after we've done this once, we're done
    quiet = 0;

    return rc;
}
