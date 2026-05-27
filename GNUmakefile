default: fmt lint test build

BINARY_BASE      := nullcloud-backend
VERSION          := 0.2.0
BINARY_VERSIONED := $(BINARY_BASE)_v$(VERSION)

# ── Platforms ──────────────────────────────────────────────────────────────────
PLATFORMS := \
	darwin_amd64 \
	darwin_arm64 \
	linux_amd64 \
	linux_arm64 \
	windows_amd64

# ── Local OS/arch for the default binary ──────────────────────────────────────
OS_NAME  := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH_RAW := $(shell uname -m)
LOCAL_OS   := $(OS_NAME)
LOCAL_ARCH := $(if $(filter arm64 aarch64,$(ARCH_RAW)),arm64,amd64)

.PHONY: default fmt lint test build clean FORCE

fmt:
	gofmt -s -w -e .

lint:
	golangci-lint run

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

## build: Cross-compile for all platforms into dist/
build: $(addprefix build-,$(PLATFORMS))

## build-<os>_<arch>: Build for one platform, e.g. make build-linux_amd64
build-%: FORCE
	$(eval GOOS   := $(word 1,$(subst _, ,$*)))
	$(eval GOARCH := $(word 2,$(subst _, ,$*)))
	$(eval EXT    := $(if $(filter windows,$(GOOS)),.exe,))
	@mkdir -p dist
	@echo "→ $(GOOS)/$(GOARCH)"
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 \
		go build -trimpath -o dist/$(BINARY_VERSIONED)_$(GOOS)_$(GOARCH)$(EXT) .

clean:
	rm -rf dist

FORCE:
