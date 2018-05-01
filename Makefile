PROJECT_NAME := Pulumi Kubernetes Resource Provider
include build/common.mk

PACK             := kubernetes
PACKDIR          := pack
PROJECT          := github.com/pulumi/pulumi-kubernetes
NODE_MODULE_NAME := @pulumi/kubernetes

PROVIDER        := pulumi-resource-${PACK}
CODEGEN         := pulumi-gen-${PACK}
VERSION         := $(shell scripts/get-version)
KUBE_VERSION    ?= v1.9.7
SWAGGER_URL     ?= https://github.com/kubernetes/kubernetes/raw/${KUBE_VERSION}/api/openapi-spec/swagger.json
OPENAPI_DIR     := pkg/gen/openapi-specs
OPENAPI_FILE    := ${OPENAPI_DIR}/swagger-${KUBE_VERSION}.json

VERSION_FLAGS   := -ldflags "-X github.com/pulumi/pulumi-kubernetes/pkg/version.Version=${VERSION}"

GO              ?= go
GOMETALINTERBIN ?= gometalinter
GOMETALINTER    :=${GOMETALINTERBIN} --config=Gometalinter.json
CURL            ?= curl

TESTPARALLELISM := 10
TESTABLE_PKGS   := ./pkg/await ./pkg/provider

$(OPENAPI_FILE)::
	@mkdir -p pkg/gen/openapi-specs
	$(CURL) -s -L $(SWAGGER_URL) > $(OPENAPI_FILE)

build:: $(OPENAPI_FILE)
	$(GO) install $(VERSION_FLAGS) $(PROJECT)/cmd/$(PROVIDER)
	$(GO) install $(VERSION_FLAGS) $(PROJECT)/cmd/$(CODEGEN)
	for LANGUAGE in "nodejs" ; do \
		$(CODEGEN) $(OPENAPI_FILE) pkg/gen/node-templates $(PACKDIR)/$$LANGUAGE || exit 3 ; \
	done

lint::
	$(GOMETALINTER) ./cmd/... resources.go | sort ; exit "$${PIPESTATUS[0]}"

test_all::
	PATH=$(PULUMI_BIN):$(PATH) go test -v -cover -timeout 1h -parallel ${TESTPARALLELISM} $(TESTABLE_PKGS)
