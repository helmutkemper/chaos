
GO ?= go
GOTESTSUM ?= gotestsum

.PHONY: build
## run the example
build:
	@$(GO) mod tidy
	@$(GO) test git_test.go -v
