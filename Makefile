SHELL=/bin/bash
.SHELLFLAGS=-e

PROJECT         = github.com/pulumi/terraform-bridge
TFGEN           = lumi-tfgen
TFGEN_BIN       = ${GOPATH}/bin/${TFGEN}
TFGEN_PKG       = ${PROJECT}/cmd/${TFGEN}
TFBRIDGE        = lumi-tfbridge
TFBRIDGE_BIN    = ${GOPATH}/bin/${TFBRIDGE}
TFBRIDGE_PKG    = ${PROJECT}/cmd/${TFBRIDGE}
GOPKGS          = $(shell go list ./cmd/... ./pkg/... | grep -v /vendor/)
LUMIROOT       ?= /usr/local/lumi
LUMILIB         = ${LUMIROOT}/packs
LUMIPLUG        = lumi-resource
TESTPARALLELISM = 10

ECHO=echo -e
GOMETALINTERBIN=gometalinter
GOMETALINTER=${GOMETALINTERBIN} --config=Gometalinter.json

all: banner tools packs
.PHONY: all

banner:
	@$(ECHO) "\033[1;37m=====================\033[0m"
	@$(ECHO) "\033[1;37mLumi Terraform Bridge\033[0m"
	@$(ECHO) "\033[1;37m=====================\033[0m"
	@go version
.PHONY: banner

$(TFGEN_BIN) tfgen:
	go install ${PROJECT}/cmd/lumi-tfgen
$(TFBRIDGE_BIN) tfbridge:
	go install ${PROJECT}/cmd/lumi-tfbridge
.PHONY: $(TFGEN_BIN) tfgen $(TFBRIDGE_BIN) tfbridge

build: $(TFGEN_BIN) $(TFBRIDGE_BIN)
.PHONY: build

tools: build test
.PHONY: tools

test:
	go test -cover -parallel ${TESTPARALLELISM} ${GOPKGS}
	which ${GOMETALINTERBIN} >/dev/null
	$(GOMETALINTER) ./cmd/... ./pkg/... | sort ; exit "$${PIPESTATUS[0]}"
	go tool vet -printf=false cmd/ pkg/
.PHONY: test

packs:
	$(MAKE) packs/aws
	$(MAKE) packs/azure
	$(MAKE) packs/gcp
.PHONY: packs

clean:
	rm -rf ${GOPATH}/bin/${TFGEN}
	rm -rf ${GOPATH}/bin/${TFBRIDGE}
.PHONY: clean

