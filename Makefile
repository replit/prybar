.PHONY: all
all: $(addprefix prybar-,$(filter-out R,$(shell ls languages))) ## Build all Prybar binaries

# Avoid depending on the subdirectories of 'languages', because their
# mtimes are updated every time we build.
prybar-%: utils/*.go languages/%/* ## Build the Prybar binary for LANG
	@echo "build prybar-$(*)"
	@if [ -f "languages/$(*)/compile" ]; then languages/$(*)/compile; fi
	@scripts/inject.sh $(*)
	@go generate languages/$(*)/main.go
	@go build -o prybar-$(*) ./languages/$(*)
	@rm -f languages/$(*)/generated_*.go

prybar-python38: utils/*.go languages/python3/*
	@echo "build prybar-python38..."
	@scripts/inject.sh python3
	CGO_CFLAGS="$(shell pkg-config --cflags python-3.8-embed) -DPYTHON_3_8" \
	CGO_LDFLAGS="$(shell pkg-config --libs python-3.8-embed)" \
	go build -o prybar-python38 ./languages/python3/
	@rm -f languages/python3/generated_*.go

prybar-python310: utils/*.go languages/python3/*
	@echo "build prybar-python310..."
	@scripts/inject.sh python3
	CGO_CFLAGS="$(shell pkg-config --cflags python-3.10-embed) -DPYTHON_3_10" \
	CGO_LDFLAGS="$(shell pkg-config --libs python-3.10-embed)" \
	go build -o prybar-python310 ./languages/python3/
	@rm -f languages/python3/generated_*.go

prybar-python311: utils/*.go languages/python3/*
	@echo "build prybar-python311..."
	@scripts/inject.sh python3
	CGO_CFLAGS="$(shell pkg-config --cflags python-3.11-embed) -DPYTHON_3_11" \
	CGO_LDFLAGS="$(shell pkg-config --libs python-3.11-embed)" \
	go build -o prybar-python311 ./languages/python3/
	@rm -f languages/python3/generated_*.go

.PHONY: docker
docker: ## Run a shell with Prybar inside Docker
	docker build . -f Dockerfile.dev -t prybar-dev
	docker run -it --rm -v "$$PWD:/gocode/src/github.com/replit/prybar" prybar-dev

.PHONY: image
image: ## Build a Docker image with Prybar for distribution
	docker build . -t prybar

.PHONY: test
test: ## Run integration tests
	./run_tests

.PHONY: test-image
test-image: image ## Test Docker image for distribution
	docker run -t --rm prybar ./run_tests

.PHONY: clean
clean: ## Remove build artifacts
	rm -f prybar-* languages/*/generated_*.go prybar_assets/sqlite/patch.so

.PHONY: help
help: ## Show this message
	@echo "usage:" >&2
	@grep -h "[#]# " $(MAKEFILE_LIST)	| \
		sed 's/^/  make /'		| \
		sed 's/:[^#]*[#]# /|/'		| \
		sed 's/%/LANG/'			| \
		column -t -s'|' >&2
