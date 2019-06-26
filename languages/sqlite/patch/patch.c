#define _GNU_SOURCE
#include <dlfcn.h>
#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// int (*orig_printf)(const char *format, ...) = NULL;
int (*orig_printf_chk)(int flag, const char *format, ...) = NULL;

// int fputc(int c, FILE *stream) {
//     if (orig_printf == NULL) {
//         orig_printf = dlsym(RTLD_NEXT, "printf");
//     }
//     return orig_printf("fputc\n");
// }
// int fputs(const char *s, FILE *stream) {
//     if (orig_printf == NULL) {
//         orig_printf = dlsym(RTLD_NEXT, "printf");
//     }
//     return orig_printf("fputs\n");
// }
// int putc(int c, FILE *stream) {
//     if (orig_printf == NULL) {
//         orig_printf = dlsym(RTLD_NEXT, "printf");
//     }
//     return orig_printf("putc\n");
// }
// int putchar(int c) {
//     if (orig_printf == NULL) {
//         orig_printf = dlsym(RTLD_NEXT, "printf");
//     }
//     return orig_printf("putchar\n");
// }

// int puts(const char *s) {
//     if (orig_printf == NULL) {
//         orig_printf = dlsym(RTLD_NEXT, "printf");
//     }
//     return orig_printf("puts\n");
// }

// int printf(const char *format, ...) {
//     if (orig_printf == NULL) {
//         orig_printf = dlsym(RTLD_NEXT, "printf");
//     }
//     return orig_printf("printf\n");
//     // if (strstr(format, "SQLite version")) {
//     // return orig_printf("hey\n");
//     // }
//     va_list arg;
//     va_start(arg, format);
//     int rc = vprintf(format, arg);
//     va_end(arg);
//     return rc;
// }

int __printf_chk(int flag, const char *format, ...) {
    if (!orig_printf_chk) {
        orig_printf_chk = dlsym(RTLD_NEXT, "__printf_chk");
    }

    char *quiet = getenv("PRYBAR_QUIET");
    if (quiet) {
        return 0;
    }

    va_list arg;
    va_start(arg, format);
    int rc = vprintf(format, arg);
    va_end(arg);
    return rc;
}

// size_t fwrite(const void *ptr, size_t size, size_t nmemb, FILE *stream) {
//     if (orig_printf == NULL) {
//         orig_printf = dlsym(RTLD_NEXT, "printf");
//     }
//     return orig_printf("fwrite\n");
// }
