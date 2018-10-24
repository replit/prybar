LANGDIR := ./languages
LANGS   := $(shell ls $(LANGDIR))
BINS    := $(addprefix prybar-,$(LANGS))

.PHONY: clean test all

all: $(BINS)

prybar-%: ./languages/$(*) ./utils/* ./linenoise/* ./languages/$(*)/*
	cp inject_launch.go ./languages/$(*)/inject_launch.go
	CGO_LDFLAGS_ALLOW=".*" go build -o prybar-$(*) ./languages/$(*)
	rm ./languages/$(*)/inject_launch.go

test:
	./run_tests

clean:
	@rm ./prybar-*
