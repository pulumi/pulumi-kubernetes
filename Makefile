MAKEFLAGS    := --warn-undefined-variables

PROJECT_NAME := Pulumi Kubernetes Resource Provider

PACK             := kubernetes
PACKDIR          := sdk
PROJECT          := github.com/pulumi/pulumi-kubernetes
NODE_MODULE_NAME := @pulumi/kubernetes
NUGET_PKG_NAME   := Pulumi.Kubernetes

PROVIDER_VERSION ?= "1.0.0-alpha.0+dev"

PROVIDER        := pulumi-resource-${PACK}
CODEGEN         := pulumi-gen-${PACK}
PROVIDER_PATH   := provider/v3
VERSION_PATH     := ${PROVIDER_PATH}/pkg/version.Version

KUBE_VERSION    ?= v1.26.0
SWAGGER_URL     ?= https://github.com/kubernetes/kubernetes/raw/${KUBE_VERSION}/api/openapi-spec/swagger.json
OPENAPI_DIR     := provider/pkg/gen/openapi-specs
OPENAPI_FILE    := ${OPENAPI_DIR}/swagger-${KUBE_VERSION}.json
SCHEMA_FILE     := provider/cmd/pulumi-resource-kubernetes/schema.json

GOPATH			:= $(shell go env GOPATH)

JAVA_GEN 		 := pulumi-java-gen
JAVA_GEN_VERSION := v0.5.2

WORKING_DIR     := $(shell pwd)
CODEGEN_PATH    = bin/${CODEGEN}
pulumictl := bin/pulumictl
TESTPARALLELISM := 4

# The general form for this Makefile is
#  - faithfully represent dependencies between build targets as files, where possible (i.e., like a normal Makefile)
#  - use phony targets to give the CI system a kind of API for building in stages

.PHONY: default
default: build

# Make sure necessary tools are present and the working dir is ready to build
.PHONY: ensure
ensure: ${pulumictl}
	cd provider && go mod tidy
	cd sdk && go mod tidy
	cd tests && go mod tidy

${pulumictl}: PULUMICTL_VERSION := $(shell cat .pulumictl.version)
${pulumictl}: PLAT := $(shell go version | sed -En "s/go version go.* (.*)\/(.*)/\1-\2/p")
${pulumictl}: PULUMICTL_URL := "https://github.com/pulumi/pulumictl/releases/download/v$(PULUMICTL_VERSION)/pulumictl-v$(PULUMICTL_VERSION)-$(PLAT).tar.gz"
${pulumictl}: .pulumictl.version
	@echo "Installing pulumictl"
	@mkdir -p bin
	wget -q -O - "$(PULUMICTL_URL)" | tar -xzf - -C $(WORKING_DIR)/bin pulumictl
	@touch ${pulumictl}
	@echo "pulumictl" $$(./bin/pulumictl version)

${OPENAPI_FILE}:
	@mkdir -p $(OPENAPI_DIR)
	curl -s -L $(SWAGGER_URL) > $(OPENAPI_FILE)

${CODEGEN_PATH}: provider/go.mod provider/cmd/$(CODEGEN)/*.go $(shell find provider/pkg -name '*.go')
	(cd provider && CGO_ENABLED=1 go build -o $(WORKING_DIR)/${CODEGEN_PATH} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${PROVIDER_VERSION}" ${PROJECT}/${PROVIDER_PATH}/cmd/$(CODEGEN))

${SCHEMA_FILE}: ${CODEGEN_PATH} ${OPENAPI_FILE}
	${CODEGEN_PATH} schema $(OPENAPI_FILE) $(CURDIR) # magically writes to the expected place

bin/${PROVIDER}: ${CODEGEN_PATH} ${SCHEMA_FILE}
	${CODEGEN_PATH} kinds $(SCHEMA_FILE) $(CURDIR) # TODO should be its own rule?
	@[ ! -f "provider/cmd/${PROVIDER}/schema.go" ] || \
		(echo "\n    Please remove provider/cmd/${PROVIDER}/schema.go, which is no longer used\n" && false)
	(cd provider && VERSION=${PROVIDER_VERSION} go generate cmd/${PROVIDER}/main.go)
	(cd provider && CGO_ENABLED=0 go build -o $(WORKING_DIR)/bin/${PROVIDER} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${PROVIDER_VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

.PHONY: k8sprovider
k8sprovider: bin/${PROVIDER}

.PHONY: k8sprovider_debug
k8sprovider_debug:
	(cd provider && CGO_ENABLED=0 go build -o $(WORKING_DIR)/bin/${PROVIDER} -gcflags="all=-N -l" -ldflags "-X ${PROJECT}/${VERSION_PATH}=${PROVIDER_VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

.PHONY: test_provider
test_provider:
	cd provider/pkg && go test -short -v -count=1 -cover -timeout 2h -parallel ${TESTPARALLELISM} ./...

.PHONY: dotnet_sdk
dotnet_sdk: DOTNET_VERSION = $(shell ${pulumictl} convert-version --language dotnet --version "${PROVIDER_VERSION}")
dotnet_sdk: ${CODEGEN_PATH} ${SCHEMA_FILE}
	${CODEGEN_PATH} -version=${DOTNET_VERSION} dotnet $(SCHEMA_FILE) $(CURDIR)
	rm -rf sdk/dotnet/bin/Debug
	cd ${PACKDIR}/dotnet/&& \
		echo "module fake_dotnet_module // Exclude this directory from Go tools\n\ngo 1.17" > go.mod && \
		echo "${DOTNET_VERSION}" >version.txt && \
		dotnet build /p:Version=${DOTNET_VERSION}

.PHONY: go_sdk
go_sdk: ${CODEGEN_PATH} ${SCHEMA_FILE}
	# Delete generated SDK before regenerating.
	rm -rf sdk/go/kubernetes
	$(WORKING_DIR)/bin/$(CODEGEN) -version=${PROVIDER_VERSION} go $(SCHEMA_FILE) $(CURDIR)

.PHONY: nodejs_sdk
nodejs_sdk: JS_VERSION = $(shell ${pulumictl} convert-version --language javascript --version "${PROVIDER_VERSION}")
nodejs_sdk: ${CODEGEN_PATH} ${SCHEMA_FILE}
	${CODEGEN_PATH} -version=${JS_VERSION} nodejs $(SCHEMA_FILE) $(CURDIR)
	cd ${PACKDIR}/nodejs/ && \
		echo "module fake_nodejs_module // Exclude this directory from Go tools\n\ngo 1.17" > go.mod && \
		yarn install && \
		yarn run tsc
	cp README.md LICENSE ${PACKDIR}/nodejs/package.json ${PACKDIR}/nodejs/yarn.lock ${PACKDIR}/nodejs/bin/
	sed -i.bak 's/$${VERSION}/$(JS_VERSION)/g' ${PACKDIR}/nodejs/bin/package.json

.PHONY: python_sdk
python_sdk: PYPI_VERSION = $(shell ${pulumictl} convert-version --language python --version "${PROVIDER_VERSION}")
python_sdk: ${CODEGEN_PATH} ${SCHEMA_FILE}
	# Delete only files and folders that are generated.
	rm -r sdk/python/pulumi_kubernetes/*/ sdk/python/pulumi_kubernetes/__init__.py
	${CODEGEN_PATH} -version=${PROVIDER_VERSION} python $(SCHEMA_FILE) $(CURDIR)
	cp README.md ${PACKDIR}/python/
	cd ${PACKDIR}/python/ && \
		echo "module fake_python_module // Exclude this directory from Go tools\n\ngo 1.17" > go.mod && \
		python3 setup.py clean --all 2>/dev/null && \
		rm -rf ./bin/ ../python.bin/ && cp -R . ../python.bin && mv ../python.bin ./bin && \
		sed -i.bak -e 's/^VERSION = .*/VERSION = "$(PYPI_VERSION)"/g' -e 's/^PLUGIN_VERSION = .*/PLUGIN_VERSION = "$(VERSION)"/g' ./bin/setup.py && \
		rm ./bin/setup.py.bak && \
		cd ./bin && python3 setup.py build sdist

.PHONY: java_sdk
java_sdk: PACKAGE_VERSION = $(shell ${pulumictl} convert-version --language generic --version "${PROVIDER_VERSION}")
java_sdk: bin/pulumi-java-gen ${CODEGEN_PATH} ${SCHEMA_FILE}
	$(WORKING_DIR)/bin/$(JAVA_GEN) generate --schema $(SCHEMA_FILE) --out sdk/java --build gradle-nexus
	cd ${PACKDIR}/java/ && \
		echo "module fake_java_module // Exclude this directory from Go tools\n\ngo 1.17" > go.mod && \
		gradle --console=plain build

bin/pulumi-java-gen:
	${pulumictl} download-binary -n pulumi-language-java -v $(JAVA_GEN_VERSION) -r pulumi/pulumi-java

.PHONY: build
build: nodejs_sdk go_sdk python_sdk dotnet_sdk java_sdk

# Required for the codegen action that runs in pulumi/pulumi
.PHONY: only_build
only_build: build

.PHONY: lint
lint:
	for DIR in "provider" "sdk" "tests" ; do \
		pushd $$DIR && golangci-lint run -c ../.golangci.yml --timeout 10m && popd ; \
	done

.PHONY: install
install: install_nodejs_sdk install_dotnet_sdk
	cp $(WORKING_DIR)/bin/${PROVIDER} ${GOPATH}/bin

GO_TEST_FAST := go test -short -v -count=1 -cover -timeout 2h -parallel ${TESTPARALLELISM}
GO_TEST 	 := go test -v -count=1 -cover -timeout 2h -parallel ${TESTPARALLELISM}

.PHONY: test_fast
test_fast:
# TODO: re-enable this test once https://github.com/pulumi/pulumi/issues/4954 is fixed.
#	./sdk/nodejs/node_modules/mocha/bin/mocha ./sdk/nodejs/bin/tests
	cd provider/pkg && $(GO_TEST_FAST) ./...
	cd tests/sdk/nodejs && $(GO_TEST_FAST) ./...
	cd tests/sdk/python && $(GO_TEST_FAST) ./...
	cd tests/sdk/dotnet && $(GO_TEST_FAST) ./...
	cd tests/sdk/go && $(GO_TEST_FAST) ./...

.PHONY: test_all
test_all:
	cd provider/pkg && $(GO_TEST) ./...
	cd tests/sdk/nodejs && $(GO_TEST) ./...
	cd tests/sdk/python && $(GO_TEST) ./...
	cd tests/sdk/dotnet && $(GO_TEST) ./...
	cd tests/sdk/go && $(GO_TEST) ./...

.PHONY: generate_schema
generate_schema: ${SCHEMA_FILE}

.PHONY: install_dotnet_sdk
install_dotnet_sdk:
	rm -rf $(WORKING_DIR)/nuget/$(NUGET_PKG_NAME).*.nupkg
	mkdir -p $(WORKING_DIR)/nuget
	find . -name '*.nupkg' -print -exec cp -p {} ${WORKING_DIR}/nuget \;

.PHONY: install_python_sdk
install_python_sdk:
	#target intentionally blank

.PHONY: install_go_sdk
install_go_sdk::
	#target intentionally blank

.PHONY: install_java_sdk
install_java_sdk:
	#target intentionally blank

.PHONY: install_nodejs_sdk
install_nodejs_sdk:
	-yarn unlink --cwd $(WORKING_DIR)/sdk/nodejs/bin
	yarn link --cwd $(WORKING_DIR)/sdk/nodejs/bin

.PHONY: examples
examples:
	cd provider/pkg/gen/examples/upstream && go run generate.go ./yaml ./
