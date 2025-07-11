APP_NAME := opml2pptx
DEFAULT_VERSION := 0.0.1-dev
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "$(DEFAULT_VERSION)")
BUILD_DIR := dist

PLATFORMS := \
	darwin/amd64 \
	darwin/arm64 \
	linux/amd64 \
	windows/amd64

.PHONY: all clean test

all: clean $(PLATFORMS)

$(PLATFORMS):
	@osarch=$@; \
	OS=$${osarch%/*}; \
	ARCH=$${osarch##*/}; \
	EXT=""; \
	[ "$$OS" = "windows" ] && EXT=".exe"; \
	OUTPUT="$(BUILD_DIR)/$$OS/$$ARCH/$(APP_NAME)$$EXT"; \
	echo "Building $$OUTPUT..."; \
	GOOS=$$OS GOARCH=$$ARCH CGO_ENABLED=0 go build -o $$OUTPUT -ldflags "-s -w -X main.version=$(VERSION)" ./cmd/$(APP_NAME)

clean:
	@rm -rf $(BUILD_DIR)

test:
	@go test ./...
