GO:=go
CMD=gitPractice
OUTPUT=./bin/$(CMD)
CURPKG := $(shell go list ./... | grep -v /vendor -m1)
PKGS := $(shell go list ./... | grep -v /vendor)

VER := $(shell git describe --abbrev=0 --tags)
REV := $(shell git rev-parse --short HEAD)
ENV := $(shell git rev-parse --abbrev-ref HEAD)

GOMETALINTER := $(BIN_DIR)/gometalinter

PHONY: build run test lint

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

build:
	@-rm $(OUTPUT)
	@echo "Building application with Version=$(VER)\tRevision=$(REV)\tEnvironment=$(ENV)"
	@$(GO) build -ldflags "-X $(CURPKG).Version=$(VER) -X $(CURPKG).Revision=$(REV) -X $(CURPKG).Environment=$(ENV)" -o $(OUTPUT) ./cmd/

run: build
	@$(OUTPUT)

test:
	@go test $(PKGS)

lint: $(GOMETALINTER)
	gometalinter ./... --vendor

