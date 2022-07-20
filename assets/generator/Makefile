BINDIR:=bin

ROOT_PACKAGE:=$(shell go list .)
COMMAND_PACKAGES:=$(shell go list ./cmd/...)

BINARIES:=$(COMMAND_PACKAGES:$(ROOT_PACKAGE)/cmd/%=$(BINDIR)/%)

GO_FILES:=$(shell find . -type f -name '*.go' -print)

ASSETS:=$(shell seq -f "../icon%g.png" 0 50)

.PHONY: assets
assets: $(ASSETS)

$(ASSETS): $(BINARIES)
	$(BINARIES) > $@

$(BINARIES): $(GO_FILES)
	@go build -o $@ $(@:$(BINDIR)/%=$(ROOT_PACKAGE)/cmd/%)
