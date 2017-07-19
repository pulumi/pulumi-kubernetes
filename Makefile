SHELL=/bin/bash
.SHELLFLAGS=-e

PROJECT         = github.com/pulumi/terraform-bridge
GOPKGS          = $(shell go list ./pkg/... | grep -v /vendor/)
TESTPARALLELISM = 10

ECHO=echo -e
GOMETALINTERBIN=gometalinter
GOMETALINTER=${GOMETALINTERBIN} --config=Gometalinter.json

all: banner build test packs
.PHONY: all

banner:
	@$(ECHO) "\033[1;37m=============================\033[0m"
	@$(ECHO) "\033[1;37mLumi Terraform Bridge\033[0m"
	@$(ECHO) "\033[1;37m=============================\033[0m"
	@go version
.PHONY: banner

build:
	go build ${PROJECT}/pkg/tfgen
	go build ${PROJECT}/pkg/tfbridge
.PHONY: build

test:
	go test -cover -parallel ${TESTPARALLELISM} ${GOPKGS}
	which ${GOMETALINTERBIN} >/dev/null
	$(GOMETALINTER) ./pkg/... | sort ; exit "$${PIPESTATUS[0]}"
	go tool vet -printf=false pkg/
.PHONY: test

packs:
	cd ./packs/aws; $(MAKE)
	cd ./packs/azure; $(MAKE)
	cd ./packs/gcp; $(MAKE)
.PHONY: packs

