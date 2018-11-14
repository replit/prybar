LANGDIR := ./languages
LANGS   := $(shell ls $(LANGDIR))
BINS    := $(addprefix prybar-,$(LANGS))

.PHONY: clean test all

all: $(BINS)

prybar-%: ./languages/$(*) ./utils/* ./linenoise/* ./languages/$(*)/*
	./scripts/inject.sh $(*)
	go generate ./languages/$(*)/main.go
	CGO_LDFLAGS_ALLOW=".*" go build -o prybar-$(*) ./languages/$(*)
	rm ./languages/$(*)/generated_*.go

test:
	./run_tests

clean:
	@rm ./prybar-* ./languages/*/generated_*.go
