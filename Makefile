VERSION ?= dev
MODULE = github.com/laiye-ai/adp-cli
LDFLAGS = -s -w -X $(MODULE)/cmd.version=$(VERSION)
DIST = dist

PLATFORMS = \
	windows/amd64/.exe \
	windows/arm64/.exe \
	linux/amd64/ \
	linux/arm64/ \
	darwin/amd64/ \
	darwin/arm64/

.PHONY: all clean build-all

all: build-all

build-all: clean
	@mkdir -p $(DIST)
	@$(foreach p,$(PLATFORMS), \
		$(eval OS := $(word 1,$(subst /, ,$(p)))) \
		$(eval ARCH := $(word 2,$(subst /, ,$(p)))) \
		$(eval EXT := $(word 3,$(subst /, ,$(p)))) \
		echo "Building $(OS)/$(ARCH)..." && \
		GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(DIST)/adp-$(OS)-$(ARCH)$(EXT) . && \
	) true

clean:
	rm -rf $(DIST)
