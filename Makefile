PROJECT_NAME := Pulumi Kubernetes Resource Provider
include build/common.mk

PACK             := kubernetes
PACKDIR          := pack
PROJECT          := github.com/pulumi/pulumi-kubernetes
NODE_MODULE_NAME := @pulumi/kubernetes

PROVIDER        := pulumi-resource-${PACK}
VERSION         := $(shell scripts/get-version)

VERSION_FLAGS   := -ldflags "-X github.com/pulumi/pulumi-kubernetes/pkg/version.Version=${VERSION}"

GOMETALINTERBIN=gometalinter
GOMETALINTER=${GOMETALINTERBIN} --config=Gometalinter.json

TESTPARALLELISM := 10

build::
	go install $(VERSION_FLAGS) ${PROJECT}/cmd/${PROVIDER}

lint::
	$(GOMETALINTER) ./cmd/... resources.go | sort ; exit "$${PIPESTATUS[0]}"
