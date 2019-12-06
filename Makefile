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

GENNY_REF := df3d48aa411e8e28498e452e42cc1e40498a9e58
bin/genny: bin/gobin
	GOBIN=${CURDIR}/bin \
	bin/gobin github.com/cheekybits/genny@$(GENNY_REF)
bins += bin/genny

bin/addgeneratedheader: gobuildcache
	cd tools; ${GOBUILD} -o ../$@ ./addgeneratedheader
bins += bin/addgeneratedheader

.PHONY: all
all: $(bins)

.PHONY: clean
clean:
	rm -rf $(bins) $(cleanup_extras)
