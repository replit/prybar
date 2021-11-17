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
