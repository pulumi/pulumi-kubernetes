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
PYPI_VERSION    := $(shell cd scripts && ./get-py-version)
KUBE_VERSION    ?= v1.18.0
SWAGGER_URL     ?= https://github.com/kubernetes/kubernetes/raw/${KUBE_VERSION}/api/openapi-spec/swagger.json
OPENAPI_DIR     := provider/pkg/gen/openapi-specs
OPENAPI_FILE    := ${OPENAPI_DIR}/swagger-${KUBE_VERSION}.json
SCHEMA_FILE     := provider/cmd/pulumi-resource-kubernetes/schema.json

VERSION_FLAGS   := -ldflags "-X github.com/pulumi/pulumi-kubernetes/provider/v2/pkg/version.Version=${VERSION}"

GO              ?= go
CURL            ?= curl
PYTHON          ?= python3

DOTNET_PREFIX  := $(firstword $(subst -, ,${VERSION:v%=%})) # e.g. 1.5.0
DOTNET_SUFFIX  := $(word 2,$(subst -, ,${VERSION:v%=%}))    # e.g. alpha.1

ifeq ($(strip ${DOTNET_SUFFIX}),)
	DOTNET_VERSION := $(strip ${DOTNET_PREFIX})
else
	DOTNET_VERSION := $(strip ${DOTNET_PREFIX})-$(strip ${DOTNET_SUFFIX})
endif

TESTPARALLELISM := 1

$(OPENAPI_FILE)::
	@mkdir -p $(OPENAPI_DIR)
	test -f $(OPENAPI_FILE) || $(CURL) -s -L $(SWAGGER_URL) > $(OPENAPI_FILE)

k8sgen::
	cd provider && $(GO) install $(VERSION_FLAGS) $(PROJECT)/provider/v2/cmd/$(CODEGEN)

$(SCHEMA_FILE):: k8sgen $(OPENAPI_FILE)
	$(call STEP_MESSAGE)
	@echo "Generating Pulumi schema..."
	$(CODEGEN) schema $(OPENAPI_FILE) $(CURDIR)
	@echo "Finished generating schema."

k8sprovider:: $(SCHEMA_FILE)
	$(CODEGEN) -version=${VERSION} kinds $(SCHEMA_FILE) $(CURDIR)
	cd provider && $(GO) generate cmd/${PROVIDER}/main.go
	cd provider && $(GO) install $(VERSION_FLAGS) $(PROJECT)/provider/v2/cmd/$(PROVIDER)

dotnet_sdk:: k8sgen $(OPENAPI_FILE)
	$(CODEGEN) -version=${VERSION} dotnet $(OPENAPI_FILE) $(CURDIR)
	cd ${PACKDIR}/dotnet/&& \
		echo "${VERSION:v%=%}" >version.txt && \
		dotnet build /p:Version=${DOTNET_VERSION}

go_sdk:: k8sgen $(SCHEMA_FILE)
	$(CODEGEN) -version=${VERSION} go $(SCHEMA_FILE) $(CURDIR)

nodejs_sdk:: k8sgen $(SCHEMA_FILE)
	$(CODEGEN) -version=${VERSION} nodejs $(SCHEMA_FILE) $(CURDIR)
	cd ${PACKDIR}/nodejs/ && \
		yarn install && \
		yarn run tsc
	cp README.md LICENSE ${PACKDIR}/nodejs/package.json ${PACKDIR}/nodejs/yarn.lock ${PACKDIR}/nodejs/bin/
	sed -i.bak 's/$${VERSION}/$(VERSION)/g' ${PACKDIR}/nodejs/bin/package.json

python_sdk:: k8sgen $(SCHEMA_FILE)
	# Delete only files and folders that are generated.
	rm -r sdk/python/pulumi_kubernetes/*/ sdk/python/pulumi_kubernetes/__init__.py
	$(CODEGEN) -version=${VERSION} python $(SCHEMA_FILE) $(CURDIR)
	cp README.md ${PACKDIR}/python/
	cd ${PACKDIR}/python/ && \
		$(PYTHON) setup.py clean --all 2>/dev/null && \
		rm -rf ./bin/ ../python.bin/ && cp -R . ../python.bin && mv ../python.bin ./bin && \
		sed -i.bak -e "s/\$${VERSION}/$(PYPI_VERSION)/g" -e "s/\$${PLUGIN_VERSION}/$(VERSION)/g" ./bin/setup.py && \
		rm ./bin/setup.py.bak && \
		cd ./bin && $(PYTHON) setup.py build sdist

.PHONY: build
build:: k8sgen k8sprovider dotnet_sdk go_sdk nodejs_sdk python_sdk

lint::
	for DIR in "provider" "sdk" "tests" ; do \
		pushd $$DIR && golangci-lint run -c ../.golangci.yml --timeout 10m && popd ; \
	done

install::
	cd provider && GOBIN=$(PULUMI_BIN) $(GO) install $(VERSION_FLAGS) $(PROJECT)/provider/v2/cmd/$(PROVIDER)
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
	cd provider/pkg && $(GO_TEST_FAST) ./...
	cd tests && $(GO_TEST_FAST) ./...

test_all::
	cd provider/pkg && $(GO_TEST_FAST) ./...
	cd tests && $(GO_TEST) ./...

generate_schema:: $(SCHEMA_FILE)
	cp $(SCHEMA_FILE) $(PACKDIR)/schema/schema.json

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
