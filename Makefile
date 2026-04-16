VERSION ?= dev
MODULE = github.com/laiye-ai/adp-cli
LDFLAGS = -s -w -X $(MODULE)/cmd.version=$(VERSION)
DIST = dist

.PHONY: all clean build-all

all: build-all

build-all: clean
	@mkdir -p $(DIST)/win32-x64 $(DIST)/win32-arm64 $(DIST)/linux-x64 $(DIST)/linux-arm64 $(DIST)/darwin-x64 $(DIST)/darwin-arm64
	@echo "Building windows/amd64..." && GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(DIST)/win32-x64/adp.exe .
	@echo "Building windows/arm64..." && GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(DIST)/win32-arm64/adp.exe .
	@echo "Building linux/amd64..."   && GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(DIST)/linux-x64/adp .
	@echo "Building linux/arm64..."   && GOOS=linux   GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(DIST)/linux-arm64/adp .
	@echo "Building darwin/amd64..."  && GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(DIST)/darwin-x64/adp .
	@echo "Building darwin/arm64..."  && GOOS=darwin  GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(DIST)/darwin-arm64/adp .

clean:
	rm -rf $(DIST)
