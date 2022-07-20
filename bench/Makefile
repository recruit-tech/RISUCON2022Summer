BINDIR:=bin

ROOT_PACKAGE:=$(shell go list .)
COMMAND_PACKAGES:=$(shell go list ./cmd/...)

BINARIES:=$(COMMAND_PACKAGES:$(ROOT_PACKAGE)/cmd/%=$(BINDIR)/%)

GO_FILES:=$(shell find . -type f -name '*.go' -print)

.PHONY: build
build: $(BINARIES)

$(BINARIES): $(GO_FILES)
	@go build -o $@ $(@:$(BINDIR)/%=$(ROOT_PACKAGE)/cmd/%)
