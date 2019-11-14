UNAME := $(shell uname | tr A-Z a-z)
GOCMD=go
GOBUILD=$(GOCMD) build

.PHONY: gobuildcache

bin/golangci-lint: bin/bindownloader
	bin/bindownloader $@
bins += bin/golangci-lint
cleanup_extras += bin/golangci-lint-*-$(UNAME)-*

bin/bindownloader:
	script/bootstrap-bindownloader
bins += bin/bindownloader

bin/gobin: bin/bindownloader
	bin/bindownloader $@
bins += bin/gobin

bin/shellcheck: bin/bindownloader
	bin/bindownloader $@
bins += bin/shellcheck

GOIMPORTS_REF := 8aaa1484dc108aa23dcf2d4a09371c0c9e280f6b
bin/goimports: bin/gobin
	GOBIN=${CURDIR}/bin \
	bin/gobin golang.org/x/tools/cmd/goimports@$(GOIMPORTS_REF)
bins += bin/goimports

.PHONY: all
all: $(bins)

.PHONY: clean
clean:
	rm -rf $(bins) $(cleanup_extras)
