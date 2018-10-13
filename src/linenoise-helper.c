#include <stdlib.h>
#include <linenoise.h>

void linenoiseCompletionStub(const char *, linenoiseCompletions *);
void linenoiseHintsStub(const char *, int *color, int *bold);


void linenoise_setup() {
	linenoiseSetCompletionCallback(linenoiseCompletionStub);
	linenoiseSetHintsCallback(linenoiseHintsStub);
}
