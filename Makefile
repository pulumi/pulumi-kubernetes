PROJECT_NAME := Pulumi Kubernetes Resource Provider
include build/common.mk

PACK             := kubernetes
PACKDIR          := sdk
PROJECT          := github.com/pulumi/pulumi-kubernetes
NODE_MODULE_NAME := @pulumi/kubernetes
NUGET_PKG_NAME   := Pulumi.Kubernetes

PROVIDER        := pulumi-resource-${PACK}
CODEGEN         := pulumi-gen-${PACK}
VERSION         ?= $(shell scripts/get-version)
PYPI_VERSION    := $(shell scripts/get-py-version)
KUBE_VERSION    ?= v1.17.0
SWAGGER_URL     ?= https://github.com/kubernetes/kubernetes/raw/${KUBE_VERSION}/api/openapi-spec/swagger.json
OPENAPI_DIR     := pkg/gen/openapi-specs
OPENAPI_FILE    := ${OPENAPI_DIR}/swagger-${KUBE_VERSION}.json

VERSION_FLAGS   := -ldflags "-X github.com/pulumi/pulumi-kubernetes/pkg/version.Version=${VERSION}"

GO              ?= go
CURL            ?= curl
PYTHON          ?= python3

DOTNET_PREFIX  := $(firstword $(subst -, ,${VERSION:v%=%})) # e.g. 1.5.0
DOTNET_SUFFIX  := $(word 2,$(subst -, ,${VERSION:v%=%}))    # e.g. alpha.1

ifeq ($(strip ${DOTNET_SUFFIX}),)
	DOTNET_VERSION := $(strip ${DOTNET_PREFIX})-preview
else
	DOTNET_VERSION := $(strip ${DOTNET_PREFIX})-preview-$(strip ${DOTNET_SUFFIX})
endif

TESTPARALLELISM := 10
TESTABLE_PKGS   := ./pkg/... ./examples/... ./tests/...

# Set NOPROXY to true to skip GOPROXY on 'ensure'
NOPROXY := false

$(OPENAPI_FILE)::
	@mkdir -p $(OPENAPI_DIR)
	test -f $(OPENAPI_FILE) || $(CURL) -s -L $(SWAGGER_URL) > $(OPENAPI_FILE)

build:: $(OPENAPI_FILE)
	$(GO) install $(VERSION_FLAGS) $(PROJECT)/cmd/$(PROVIDER)
	$(GO) install $(VERSION_FLAGS) $(PROJECT)/cmd/$(CODEGEN)
	# Delete only files and folders that are generated.
	rm -r sdk/python/pulumi_kubernetes/*/ sdk/python/pulumi_kubernetes/__init__.py
	for LANGUAGE in "dotnet" "nodejs" "python" ; do \
		$(CODEGEN) $$LANGUAGE $(OPENAPI_FILE) pkg/gen/$${LANGUAGE}-templates $(PACKDIR) || exit 3 ; \
	done
	cd ${PACKDIR}/nodejs/ && \
		yarn install && \
		yarn run tsc
	cp README.md LICENSE ${PACKDIR}/nodejs/package.json ${PACKDIR}/nodejs/yarn.lock ${PACKDIR}/nodejs/bin/
	cp README.md ${PACKDIR}/python/
	sed -i.bak 's/$${VERSION}/$(VERSION)/g' ${PACKDIR}/nodejs/bin/package.json
	cd ${PACKDIR}/python/ && \
		$(PYTHON) setup.py clean --all 2>/dev/null && \
		rm -rf ./bin/ ../python.bin/ && cp -R . ../python.bin && mv ../python.bin ./bin && \
		sed -i.bak -e "s/\$${VERSION}/$(PYPI_VERSION)/g" -e "s/\$${PLUGIN_VERSION}/$(VERSION)/g" ./bin/setup.py && \
		rm ./bin/setup.py.bak && \
		cd ./bin && $(PYTHON) setup.py build sdist
	cd ${PACKDIR}/dotnet/&& \
		echo "${VERSION:v%=%}" >version.txt && \
		dotnet build /p:Version=${DOTNET_VERSION}

lint::
	golangci-lint run --timeout 20m

install::
	GOBIN=$(PULUMI_BIN) $(GO) install $(VERSION_FLAGS) $(PROJECT)/cmd/$(PROVIDER)
	[ ! -e "$(PULUMI_NODE_MODULES)/$(NODE_MODULE_NAME)" ] || rm -rf "$(PULUMI_NODE_MODULES)/$(NODE_MODULE_NAME)"
	mkdir -p "$(PULUMI_NODE_MODULES)/$(NODE_MODULE_NAME)"
	cp -r sdk/nodejs/bin/. "$(PULUMI_NODE_MODULES)/$(NODE_MODULE_NAME)"
	rm -rf "$(PULUMI_NODE_MODULES)/$(NODE_MODULE_NAME)/node_modules"
	rm -rf "$(PULUMI_NODE_MODULES)/$(NODE_MODULE_NAME)/tests"
	cd "$(PULUMI_NODE_MODULES)/$(NODE_MODULE_NAME)" && \
		yarn install --offline --production && \
		(yarn unlink > /dev/null 2>&1 || true) && \
		yarn link
	echo "Copying ${NUGET_PKG_NAME} NuGet packages to ${PULUMI_NUGET}"
	mkdir -p $(PULUMI_NUGET)
	rm -rf "$(PULUMI_NUGET)/$(NUGET_PKG_NAME).*.nupkg"
	find . -name '$(NUGET_PKG_NAME).*.nupkg' -exec cp -p {} ${PULUMI_NUGET} \;

test_fast::
	./sdk/nodejs/node_modules/mocha/bin/mocha ./sdk/nodejs/bin/tests
	$(GO_TEST_FAST) $(TESTABLE_PKGS)

test_all::
	$(GO_TEST) $(TESTABLE_PKGS)

.PHONY: publish_tgz
publish_tgz:
	$(call STEP_MESSAGE)
	./scripts/publish_tgz.sh

# While pulumi-kubernetes is not built using tfgen/tfbridge, the layout of the source tree is the same as these
# style of repositories, so we can re-use the common publishing scripts.
.PHONY: publish_packages
publish_packages:
	$(call STEP_MESSAGE)
	$$(go env GOPATH)/src/github.com/pulumi/scripts/ci/publish-tfgen-package .
	$$(go env GOPATH)/src/github.com/pulumi/scripts/ci/build-package-docs.sh kubernetes

.PHONY: check_clean_worktree
check_clean_worktree:
	$$(go env GOPATH)/src/github.com/pulumi/scripts/ci/check-worktree-is-clean.sh

# The travis_* targets are entrypoints for CI.
.PHONY: travis_cron travis_push travis_pull_request travis_api
travis_cron: all
travis_push: only_build check_clean_worktree publish_tgz only_test publish_packages
travis_pull_request: only_build check_clean_worktree only_test
travis_api: all
