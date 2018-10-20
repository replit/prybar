LANGDIR := ./languages
LANGS   := $(shell ls $(LANGDIR))
BINS    := $(addprefix prybar-,$(LANGS))

.PHONY: clean default

default: $(BINS)

prybar-%: ./languages/$(*) ./utils/* ./linenoise/* ./languages/$(*)/*
	cp inject_launch.go ./languages/$(*)/inject_launch.go
	CGO_LDFLAGS_ALLOW=".*" go build -o prybar-$(*) ./languages/$(*)
	rm ./languages/$(*)/inject_launch.go

clean:
	@rm ./prybar-*
