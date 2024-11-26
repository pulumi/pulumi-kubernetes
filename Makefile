PROJECT_NAME := Pulumi Kubernetes Resource Provider

PACK             := kubernetes
PACKDIR          := sdk
PROJECT          := github.com/pulumi/pulumi-kubernetes
NODE_MODULE_NAME := @pulumi/kubernetes
NUGET_PKG_NAME   := Pulumi.Kubernetes

PROVIDER        := pulumi-resource-${PACK}
CODEGEN         := pulumi-gen-${PACK}
PROVIDER_PATH   := v4/provider
VERSION_PATH     := ${PROVIDER_PATH}/pkg/version.Version

KUBE_VERSION    ?= v1.31.0
SWAGGER_URL     ?= https://github.com/kubernetes/kubernetes/raw/${KUBE_VERSION}/api/openapi-spec/swagger.json
OPENAPI_DIR     := provider/pkg/gen/openapi-specs
OPENAPI_FILE    := ${OPENAPI_DIR}/swagger-${KUBE_VERSION}.json
SCHEMA_FILE     := provider/cmd/pulumi-resource-kubernetes/schema.json
GOPATH			:= $(shell go env GOPATH)

JAVA_GEN		 := pulumi-java-gen
JAVA_GEN_VERSION := v0.13.0

WORKING_DIR     := $(shell pwd)

# Override during CI using `make [TARGET] PROVIDER_VERSION=""` or by setting a PROVIDER_VERSION environment variable
# Local & branch builds will just used this fixed default version unless specified
PROVIDER_VERSION ?= 4.0.0-alpha.0+dev
# Use this normalised version everywhere rather than the raw input to ensure consistency.
VERSION_GENERIC = $(shell pulumictl convert-version --language generic --version "$(PROVIDER_VERSION)")

openapi_file::
	@mkdir -p $(OPENAPI_DIR)
	test -f $(OPENAPI_FILE) || curl -s -L $(SWAGGER_URL) > $(OPENAPI_FILE)

ensure::
	go mod tidy

k8sgen::
	(cd provider && CGO_ENABLED=1 go build -o $(WORKING_DIR)/bin/${CODEGEN} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION_GENERIC}" ${PROJECT}/${PROVIDER_PATH}/cmd/$(CODEGEN))

schema:: k8sgen
	@echo "Generating Pulumi schema..."
	$(WORKING_DIR)/bin/${CODEGEN} schema $(OPENAPI_FILE) $(CURDIR)
	@echo "Finished generating schema."

k8sprovider::
	$(WORKING_DIR)/bin/${CODEGEN} kinds $(SCHEMA_FILE) $(CURDIR)
	@[ ! -f "provider/cmd/${PROVIDER}/schema.go" ] || \
		(echo "\n    Please remove provider/cmd/${PROVIDER}/schema.go, which is no longer used\n" && false)
	(cd provider && VERSION=${VERSION_GENERIC} go generate cmd/${PROVIDER}/main.go)
	(cd provider && CGO_ENABLED=0 go build -o $(WORKING_DIR)/bin/${PROVIDER} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION_GENERIC}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

k8sprovider_debug::
	$(WORKING_DIR)/bin/${CODEGEN} kinds $(SCHEMA_FILE) $(CURDIR)
	@[ ! -f "provider/cmd/${PROVIDER}/schema.go" ] || \
		(echo "\n    Please remove provider/cmd/${PROVIDER}/schema.go, which is no longer used\n" && false)
	(cd provider && VERSION=${VERSION_GENERIC} go generate cmd/${PROVIDER}/main.go)
	(cd provider && CGO_ENABLED=0 go build -o $(WORKING_DIR)/bin/${PROVIDER} -gcflags="all=-N -l" -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION_GENERIC}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

test_provider::
	cd provider/pkg && go test -short -v -coverprofile="coverage.txt" -coverpkg=./... -timeout 2h ./...

dotnet_sdk:: DOTNET_VERSION := $(shell pulumictl convert-version --language dotnet -v "$(VERSION_GENERIC)")
dotnet_sdk::
	$(WORKING_DIR)/bin/$(CODEGEN) -version=${VERSION_GENERIC} dotnet $(SCHEMA_FILE) $(CURDIR)
	rm -rf sdk/dotnet/bin/Debug
	cd ${PACKDIR}/dotnet/&& \
		echo "module fake_dotnet_module // Exclude this directory from Go tools\n\ngo 1.17" > go.mod && \
		dotnet build

go_sdk::
	# Delete generated SDK before regenerating.
	rm -rf sdk/go/kubernetes
	$(WORKING_DIR)/bin/$(CODEGEN) -version=${VERSION_GENERIC} go $(SCHEMA_FILE) $(CURDIR)

nodejs_sdk:: NODE_VERSION := $(shell pulumictl convert-version --language javascript -v "$(VERSION_GENERIC)")
nodejs_sdk::
	$(WORKING_DIR)/bin/$(CODEGEN) -version=${VERSION_GENERIC} nodejs $(SCHEMA_FILE) $(CURDIR)
	cd ${PACKDIR}/nodejs/ && \
		echo "module fake_nodejs_module // Exclude this directory from Go tools\n\ngo 1.17" > go.mod && \
		yarn install && \
		yarn run tsc
	cp README.md LICENSE ${PACKDIR}/nodejs/package.json ${PACKDIR}/nodejs/yarn.lock ${PACKDIR}/nodejs/bin/

python_sdk:: PYPI_VERSION := $(shell pulumictl convert-version --language python -v "$(VERSION_GENERIC)")
python_sdk::
	# Delete only files and folders that are generated.
	rm -rf sdk/python/pulumi_kubernetes/*/ sdk/python/pulumi_kubernetes/__init__.py
	# Delete files not tracked in Git
	cd ${PACKDIR}/python/ && git clean -fxd
	$(WORKING_DIR)/bin/$(CODEGEN) -version=${VERSION_GENERIC} python $(SCHEMA_FILE) $(CURDIR)
	cp README.md ${PACKDIR}/python/
	PYPI_VERSION=$(PYPI_VERSION) ./scripts/build_python_sdk.sh

java_sdk:: PACKAGE_VERSION := $(shell pulumictl convert-version --language generic -v "$(VERSION_GENERIC)")
java_sdk:: bin/pulumi-java-gen
	$(WORKING_DIR)/bin/$(JAVA_GEN) generate --schema $(SCHEMA_FILE) --overlay provider/pkg/gen/java-templates \
		--out sdk/java --build gradle-nexus
	cd ${PACKDIR}/java/ && \
		echo "module fake_java_module // Exclude this directory from Go tools\n\ngo 1.17" > go.mod && \
		gradle --console=plain build

bin/pulumi-java-gen::
	$(shell pulumictl download-binary -n pulumi-language-java -v $(JAVA_GEN_VERSION) -r pulumi/pulumi-java)

.PHONY: build
build:: k8sgen openapi_file schema k8sprovider nodejs_sdk go_sdk python_sdk dotnet_sdk java_sdk

# Required for the codegen action that runs in pulumi/pulumi
only_build:: build

lint::
	@for DIR in "provider" "tests" ; do \
		pushd $$DIR  > /dev/null; golangci-lint run -c ../.golangci.yml --timeout 10m; popd  > /dev/null; \
	done

install_provider:: k8sprovider
	cp $(WORKING_DIR)/bin/${PROVIDER} ${GOPATH}/bin

install:: install_nodejs_sdk install_dotnet_sdk install_provider

GO_TEST_FAST := go test -short -v -cover -timeout 2h
GO_TEST		 := go test -v -cover -timeout 2h

# Required for the codegen action that runs in pulumi/pulumi
test:: test_all

test_fast::
# TODO: re-enable this test once https://github.com/pulumi/pulumi/issues/4954 is fixed.
#	./sdk/nodejs/node_modules/mocha/bin/mocha ./sdk/nodejs/bin/tests
	cd provider/pkg && $(GO_TEST_FAST) ./...
	cd tests/sdk/nodejs && $(GO_TEST_FAST) ./...
	cd tests/sdk/python && $(GO_TEST_FAST) ./...
	cd tests/sdk/dotnet && $(GO_TEST_FAST) ./...
	cd tests/sdk/go && $(GO_TEST_FAST) ./...
	cd tests/sdk/java && $(GO_TEST_FAST) ./...

test_all::
	cd provider/pkg && $(GO_TEST) ./...
	cd tests/sdk/nodejs && $(GO_TEST) ./...
	cd tests/sdk/python && $(GO_TEST) ./...
	cd tests/sdk/dotnet && $(GO_TEST) ./...
	cd tests/sdk/go && $(GO_TEST) ./...
	cd tests/sdk/java && $(GO_TEST) ./...

generate_schema:: schema

install_dotnet_sdk::
	rm -rf $(WORKING_DIR)/nuget/$(NUGET_PKG_NAME).*.nupkg
	mkdir -p $(WORKING_DIR)/nuget
	find . -name '*.nupkg' -print -exec cp -p {} ${WORKING_DIR}/nuget \;

install_python_sdk::
	#target intentionally blank

install_go_sdk::
	#target intentionally blank

install_java_sdk::
	#target intentionally blank

install_nodejs_sdk::
	-yarn unlink --cwd $(WORKING_DIR)/sdk/nodejs/bin
	yarn link --cwd $(WORKING_DIR)/sdk/nodejs/bin

examples::
	cd provider/pkg/gen/examples/upstream && go run generate.go ./yaml ./
