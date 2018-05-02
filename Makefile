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
TESTABLE_PKGS   := ./pkg/await ./pkg/provider ./examples

$(OPENAPI_FILE)::
	@mkdir -p $(OPENAPI_DIR)
	$(CURL) -s -L $(SWAGGER_URL) > $(OPENAPI_FILE)

build:: $(OPENAPI_FILE)
	$(GO) install $(VERSION_FLAGS) $(PROJECT)/cmd/$(PROVIDER)
	$(GO) install $(VERSION_FLAGS) $(PROJECT)/cmd/$(CODEGEN)
	for LANGUAGE in "nodejs" ; do \
		$(CODEGEN) $(OPENAPI_FILE) pkg/gen/node-templates $(PACKDIR)/$$LANGUAGE || exit 3 ; \
	done
	cd ${PACKDIR}/nodejs/ && \
		yarn install && \
		yarn link @pulumi/pulumi && \
		yarn run tsc
	cp README.md LICENSE ${PACKDIR}/nodejs/package.json ${PACKDIR}/nodejs/yarn.lock ${PACKDIR}/nodejs/bin/

lint::
	$(GOMETALINTER) ./cmd/... ./pkg/... | sort ; exit "$${PIPESTATUS[0]}"

install::
	GOBIN=$(PULUMI_BIN) $(GO) install $(VERSION_FLAGS) $(PROJECT)/cmd/$(PROVIDER)
	[ ! -e "$(PULUMI_NODE_MODULES)/$(NODE_MODULE_NAME)" ] || rm -rf "$(PULUMI_NODE_MODULES)/$(NODE_MODULE_NAME)"
	mkdir -p "$(PULUMI_NODE_MODULES)/$(NODE_MODULE_NAME)"
	cp -r pack/nodejs/bin/. "$(PULUMI_NODE_MODULES)/$(NODE_MODULE_NAME)"
	rm -rf "$(PULUMI_NODE_MODULES)/$(NODE_MODULE_NAME)/node_modules"
	cd "$(PULUMI_NODE_MODULES)/$(NODE_MODULE_NAME)" && \
		yarn install --offline --production && \
		(yarn unlink > /dev/null 2>&1 || true) && \
		yarn link

test_all::
	PATH=$(PULUMI_BIN):$(PATH) $(GO) test -v -cover -timeout 1h -parallel ${TESTPARALLELISM} $(TESTABLE_PKGS)

# The travis_* targets are entrypoints for CI.
.PHONY: travis_cron travis_push travis_pull_request travis_api
travis_cron: all
travis_push: only_build publish_tgz only_test publish_packages
travis_pull_request: all
travis_api: all
