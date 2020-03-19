PKGS := $(shell go list ./... | grep -v /vendor)

.PHONY: test
test: lint
	go test $(PKGS)

BIN_DIR := $(GOROOT)/bin
GOLANGCILINT := $(BIN_DIR)/golangci-lint
GOIMPORTS := $(BIN_DIR)/goimports

$(GOLANGCILINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN_DIR) v1.23.8

$(GOIMPORTS):
	go get -u golang.org/x/tools/cmd/goimports

.PHONY: lint
lint: $(GOLANGCILINT)
lint: fmt
	golangci-lint run --skip-dirs-use-default

.PHONY: fmt
fmt: $(GOIMPORTS)
	goimports -w .

BINARY := afore
VERSION ?= $(shell git rev-parse --short=8 HEAD)
PLATFORMS := windows linux darwin
os = $(word 1, $@)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
	GOOS=$(os) GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-$(os)-amd64

.PHONY: release
release: test
release: windows linux darwin
