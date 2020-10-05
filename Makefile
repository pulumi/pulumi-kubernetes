PROJECT_NAME := Pulumi Kubernetes Resource Provider
include build/common.mk

PACK             := kubernetes
PACKDIR          := sdk
PROJECT          := github.com/pulumi/pulumi-kubernetes
NODE_MODULE_NAME := @pulumi/kubernetes
NUGET_PKG_NAME   := Pulumi.Kubernetes

PROVIDER        := pulumi-resource-${PACK}
CODEGEN         := pulumi-gen-${PACK}
VERSION         ?= $(shell pulumictl get version)
KUBE_VERSION    ?= v1.19.0
SWAGGER_URL     ?= https://github.com/kubernetes/kubernetes/raw/${KUBE_VERSION}/api/openapi-spec/swagger.json
OPENAPI_DIR     := provider/pkg/gen/openapi-specs
OPENAPI_FILE    := ${OPENAPI_DIR}/swagger-${KUBE_VERSION}.json
SCHEMA_FILE     := provider/cmd/pulumi-resource-kubernetes/schema.json

VERSION_FLAGS   := -ldflags "-X github.com/pulumi/pulumi-kubernetes/provider/v2/pkg/version.Version=${VERSION}"

GO              ?= go
CURL            ?= curl
PYTHON          ?= python3

WORKING_DIR     := $(shell pwd)
TESTPARALLELISM := 4

openapi_file::
	@mkdir -p $(OPENAPI_DIR)
	test -f $(OPENAPI_FILE) || $(CURL) -s -L $(SWAGGER_URL) > $(OPENAPI_FILE)

k8sgen::
	(cd provider && go build -a -o $(WORKING_DIR)/bin/${CODEGEN} $(VERSION_FLAGS) $(PROJECT)/provider/v2/cmd/$(CODEGEN))

schema::
	$(call STEP_MESSAGE)
	@echo "Generating Pulumi schema..."
	$(WORKING_DIR)/bin/${CODEGEN} schema $(OPENAPI_FILE) $(CURDIR)
	@echo "Finished generating schema."

k8sprovider::
	$(WORKING_DIR)/bin/${CODEGEN} kinds $(SCHEMA_FILE) $(CURDIR)
	(cd provider && VERSION=${VERSION} go generate cmd/${PROVIDER}/main.go)
	(cd provider && go build -a -o $(WORKING_DIR)/bin/${PROVIDER} $(VERSION_FLAGS) $(PROJECT)/provider/v2/cmd/$(PROVIDER))

dotnet_sdk:: VERSION := $(shell pulumictl get version --language dotnet)
dotnet_sdk::
	$(WORKING_DIR)/bin/$(CODEGEN) -version=${VERSION} dotnet $(SCHEMA_FILE) $(CURDIR)
	rm -rf sdk/dotnet/bin/Debug
	cd ${PACKDIR}/dotnet/&& \
		echo "${VERSION:v%=%}" >version.txt && \
		dotnet build /p:Version=${DOTNET_VERSION}

go_sdk::
	$(WORKING_DIR)/bin/$(CODEGEN) -version=${VERSION} go $(SCHEMA_FILE) $(CURDIR)

nodejs_sdk:: VERSION := $(shell pulumictl get version --language javascript)
nodejs_sdk::
	$(WORKING_DIR)/bin/$(CODEGEN) -version=${VERSION} nodejs $(SCHEMA_FILE) $(CURDIR)
	cd ${PACKDIR}/nodejs/ && \
		yarn install && \
		yarn run tsc
	cp README.md LICENSE ${PACKDIR}/nodejs/package.json ${PACKDIR}/nodejs/yarn.lock ${PACKDIR}/nodejs/bin/
	sed -i.bak 's/$${VERSION}/$(VERSION)/g' ${PACKDIR}/nodejs/bin/package.json

nodejs_sdk:: PYPI_VERSION := $(shell pulumictl get version --language python)
python_sdk::
	# Delete only files and folders that are generated.
	rm -r sdk/python/pulumi_kubernetes/*/ sdk/python/pulumi_kubernetes/__init__.py
	$(WORKING_DIR)/bin/$(CODEGEN) -version=${VERSION} python $(SCHEMA_FILE) $(CURDIR)
	cp README.md ${PACKDIR}/python/
	cd ${PACKDIR}/python/ && \
		$(PYTHON) setup.py clean --all 2>/dev/null && \
		rm -rf ./bin/ ../python.bin/ && cp -R . ../python.bin && mv ../python.bin ./bin && \
		sed -i.bak -e "s/\$${VERSION}/$(PYPI_VERSION)/g" -e "s/\$${PLUGIN_VERSION}/$(VERSION)/g" ./bin/setup.py && \
		rm ./bin/setup.py.bak && \
		cd ./bin && $(PYTHON) setup.py build sdist

.PHONY: build
build:: k8sgen openapi_file schema k8sprovider dotnet_sdk go_sdk nodejs_sdk python_sdk

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
	rm -rf $(PULUMI_NUGET)/$(NUGET_PKG_NAME).*.nupkg
	find . -name '$(NUGET_PKG_NAME).*.nupkg' -exec cp -p {} ${PULUMI_NUGET} \;

test_fast::
# TODO: re-enable this test once https://github.com/pulumi/pulumi/issues/4954 is fixed.
#	./sdk/nodejs/node_modules/mocha/bin/mocha ./sdk/nodejs/bin/tests
	cd provider/pkg && $(GO_TEST_FAST) ./...
	cd tests/sdk/nodejs && $(GO_TEST_FAST) ./...
	cd tests/sdk/python && $(GO_TEST_FAST) ./...
	cd tests/sdk/dotnet && $(GO_TEST_FAST) ./...
# TODO: re-enable Go SDK tests once CI OOM errors are fixed.
	#cd tests/sdk/go && $(GO_TEST_FAST) ./...

test_all::
	cd provider/pkg && $(GO_TEST) ./...
	cd provider/pkg && $(GO_TEST) ./...
	cd tests/sdk/nodejs && $(GO_TEST) ./...
	cd tests/sdk/python && $(GO_TEST) ./...
	cd tests/sdk/dotnet && $(GO_TEST) ./...
# TODO: re-enable Go SDK tests once CI OOM errors are fixed.
	#cd tests/sdk/go && $(GO_TEST) ./...

generate_schema:: $(SCHEMA_FILE)
