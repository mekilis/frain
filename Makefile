# Makefile
# Frain
PLATFORMS := darwin linux windows
VERSION := $(shell cat VERSION)
os = $(word 1, $@)

# Go
BINARY := frain
BUILDFLAGS := -ldflags "-X main.buildVersion=$(VERSION)"
GOBUILDDIR := cmd/frain/
GOINSTALL := $(GOCMD) install $(BUILDFLAGS)
GOFILES := $(GOBUILDDIR)*.go
GOCMD := go
GOBIN := $(GOPATH)/bin
GOBUILD := $(GOCMD) build $(BUILDFLAGS)
GOCLEAN := $(GOCMD) clean
GOGET := $(GOCMD) get
GOINSTALL := $(GOCMD) install $(BUILDFLAGS)
GORUN := $(GOCMD) run $(BUILDFLAGS)
GOTEST := $(GOCMD) test

all: test build

build:
	$(GOBUILD) -v -o $(BINARY) $(GOFILES)

clean:
	$(GOCLEAN)
	rm -f $(BINARY)
	rm -rf release

install:
	cd $(GOBUILDDIR) && $(GOINSTALL)

run:
	$(GORUN) -v $(GOFILES) -h

test:
	$(GOTEST) -v

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
	GOOS=$(os) GOARCH=amd64 $(GOBUILD) -o release/$(BINARY)-v$(VERSION)-$(os)-amd64 $(GOFILES)

.PHONY: release
release: darwin linux windows