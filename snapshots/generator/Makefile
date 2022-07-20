BINDIR:=bin

ROOT_PACKAGE:=$(shell go list .)
COMMAND_PACKAGES:=$(shell go list ./cmd/...)

BINARIES:=$(COMMAND_PACKAGES:$(ROOT_PACKAGE)/cmd/%=$(BINDIR)/%)

GO_FILES:=$(shell find . ../../bench ../../webapp/go -type f -name '*.go' -print)

SNAPSHOTS:=$(shell seq -f "../snapshot%g.json" 0 9)

.PHONY: snapshots
snapshots: $(SNAPSHOTS)

$(SNAPSHOTS): $(BINARIES)
	sleep 1
	$(BINARIES) -minify > $@

$(BINARIES): $(GO_FILES)
	@go build -o $@ $(@:$(BINDIR)/%=$(ROOT_PACKAGE)/cmd/%)
