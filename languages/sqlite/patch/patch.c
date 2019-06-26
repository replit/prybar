#define _GNU_SOURCE
#include <dlfcn.h>
#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int (*orig_printf_chk)(int flag, const char *format, ...) = NULL;
int (*orig_process_input)(void *state, FILE *in) = NULL;

int quiet;

int __printf_chk(int flag, const char *format, ...) {
    if (!orig_printf_chk) {
        orig_printf_chk = dlsym(RTLD_NEXT, "__printf_chk");

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
    int rc = vprintf(format, arg);
    va_end(arg);

    // after we've done this once, we're done
    quiet = 0;

    return rc;
}

// the first time this is called, we know that we should no longer suppress any
// output
// int process_input(void *state, FILE *in) {
//     printf("hey\n");
//     if (!orig_process_input) {
//         orig_process_input = dlsym(RTLD_NEXT, "process_input");
//     }

//     passthrough = 1;

//     return orig_process_input(state, in);
// }
