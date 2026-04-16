VERSION ?= dev
MODULE = github.com/laiye-ai/adp-cli
LDFLAGS = -s -w -X $(MODULE)/cmd.version=$(VERSION)
DIST = dist

.PHONY: all clean build-all

all: build-all

build-all: clean
	@mkdir -p $(DIST)
	@echo "Building windows/amd64..." && GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(DIST)/adp-win32-x64.exe .
	@echo "Building windows/arm64..." && GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(DIST)/adp-win32-arm64.exe .
	@echo "Building linux/amd64..."   && GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(DIST)/adp-linux-x64 .
	@echo "Building linux/arm64..."   && GOOS=linux   GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(DIST)/adp-linux-arm64 .
	@echo "Building darwin/amd64..."  && GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(DIST)/adp-darwin-x64 .
	@echo "Building darwin/arm64..."  && GOOS=darwin  GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(DIST)/adp-darwin-arm64 .

clean:
	rm -rf $(DIST)
