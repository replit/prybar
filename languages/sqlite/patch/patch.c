#define _GNU_SOURCE
#include <dlfcn.h>
#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int (*orig_fprintf_chk)(FILE *out, int flag, const char *format, ...) = NULL;

int env_checked = 0;
int quiet;

int __printf_chk(int flag, const char *format, ...) {
    if (prybar_quiet()) {
        return 0;
    }

    va_list arg;
    va_start(arg, format);
    int rc = internal_print(stdout, format, arg);
    va_end(arg);
    return rc;
}

int __fprintf_chk(FILE *out, int flag, const char *format, ...) {
    if (prybar_quiet()) {
        return 0;
    }

    va_list arg;
    va_start(arg, format);
    int rc = internal_print(out, format, arg);
    va_end(arg);

    // after we've done this once, we're done
    quiet = 0;

    return rc;
}

int internal_print(FILE *out, const char *format, va_list args) {
    return vfprintf(out, format, args);
}

int prybar_quiet() {
    if (env_checked)
        return quiet;

    if (getenv("PRYBAR_QUIET"))
        quiet = 1;
    else
        quiet = 0;

    env_checked = 1;
    return quiet;
}
