PROJECT_NAME := Pulumi Kubernetes Resource Provider

PACK             := kubernetes
PACKDIR          := sdk
PROJECT          := github.com/pulumi/pulumi-kubernetes
NODE_MODULE_NAME := @pulumi/kubernetes
NUGET_PKG_NAME   := Pulumi.Kubernetes

PROVIDER        := pulumi-resource-${PACK}
CODEGEN         := pulumi-gen-${PACK}
PROVIDER_PATH   := provider/v4
VERSION_PATH     := ${PROVIDER_PATH}/pkg/version.Version

KUBE_VERSION    ?= v1.35.1
SWAGGER_URL     ?= https://github.com/kubernetes/kubernetes/raw/${KUBE_VERSION}/api/openapi-spec/swagger.json
OPENAPI_DIR     := provider/pkg/gen/openapi-specs
OPENAPI_FILE    := ${OPENAPI_DIR}/swagger-${KUBE_VERSION}.json
SCHEMA_FILE     := provider/cmd/pulumi-resource-kubernetes/schema.json
GOPATH			:= $(shell go env GOPATH)

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
	cd provider && go mod tidy
	cd sdk && go mod tidy
	cd tests && go mod tidy

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
java_sdk::
	pulumi package gen-sdk --language java $(SCHEMA_FILE) \
		--overlays provider/pkg/gen/java-templates --out sdk
	cd ${PACKDIR}/java/ && \
		echo "module fake_java_module // Exclude this directory from Go tools\n\ngo 1.17" > go.mod && \
		gradle --console=plain build

.PHONY: build
build:: k8sgen openapi_file schema k8sprovider nodejs_sdk go_sdk python_sdk dotnet_sdk java_sdk

# Required for the codegen action that runs in pulumi/pulumi
only_build:: build

lint::
	git grep -l 'go:embed' -- provider | xargs perl -i -pe 's/go:embed/ goembed/g'
	cd provider && golangci-lint run -c ../.golangci.yml --timeout 10m
	git grep -l ' goembed' -- provider | xargs perl -i -pe 's/ goembed/go:embed/g'
	cd tests && golangci-lint run -c ../.golangci.yml --timeout 10m

# lint-fix is a utility target meant to be run manually
# that will run the linter and fix errors when possible.
lint.fix::
	git grep -l 'go:embed' -- provider | xargs perl -i -pe 's/go:embed/ goembed/g'
	cd provider && golangci-lint run -c ../.golangci.yml --fix --timeout 10m
	git grep -l ' goembed' -- provider | xargs perl -i -pe 's/ goembed/go:embed/g'
	cd tests && golangci-lint run -c ../.golangci.yml --fix --timeout 10m


install_provider:: k8sprovider
	cp $(WORKING_DIR)/bin/${PROVIDER} ${GOPATH}/bin/

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
	mkdir -p nuget
	find sdk/dotnet/bin -name '*.nupkg' -print -exec cp -p "{}" ${WORKING_DIR}/nuget \;
	if ! dotnet nuget list source | grep "${WORKING_DIR}/nuget"; then \
		dotnet nuget add source "${WORKING_DIR}/nuget" --name "${WORKING_DIR}/nuget" \
	; fi

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

# Set these variables to enable signing of the windows binary
AZURE_SIGNING_CLIENT_ID ?=
AZURE_SIGNING_CLIENT_SECRET ?=
AZURE_SIGNING_TENANT_ID ?=
AZURE_SIGNING_KEY_VAULT_URI ?=
SKIP_SIGNING ?=

bin/jsign-6.0.jar:
	wget https://github.com/ebourg/jsign/releases/download/6.0/jsign-6.0.jar --output-document=bin/jsign-6.0.jar

sign-goreleaser-exe-amd64: GORELEASER_ARCH := amd64_v1
sign-goreleaser-exe-arm64: GORELEASER_ARCH := arm64

# Set the shell to bash to allow for the use of bash syntax.
sign-goreleaser-exe-%: SHELL:=/bin/bash
sign-goreleaser-exe-%: bin/jsign-6.0.jar
	@# Only sign windows binary if fully configured.
	@# Test variables set by joining with | between and looking for || showing at least one variable is empty.
	@# Move the binary to a temporary location and sign it there to avoid the target being up-to-date if signing fails.
	@set -e; \
	if [[ "${SKIP_SIGNING}" != "true" ]]; then \
		if [[ "|${AZURE_SIGNING_CLIENT_ID}|${AZURE_SIGNING_CLIENT_SECRET}|${AZURE_SIGNING_TENANT_ID}|${AZURE_SIGNING_KEY_VAULT_URI}|" == *"||"* ]]; then \
			echo "Can't sign windows binaries as required configuration not set: AZURE_SIGNING_CLIENT_ID, AZURE_SIGNING_CLIENT_SECRET, AZURE_SIGNING_TENANT_ID, AZURE_SIGNING_KEY_VAULT_URI"; \
			echo "To rebuild with signing delete the unsigned windows exe file and rebuild with the fixed configuration"; \
			if [[ "${CI}" == "true" ]]; then exit 1; fi; \
		else \
			file=dist/build-provider-sign-windows_windows_${GORELEASER_ARCH}/pulumi-resource-kubernetes.exe; \
			mv $${file} $${file}.unsigned; \
			az login --service-principal \
				--username "${AZURE_SIGNING_CLIENT_ID}" \
				--password "${AZURE_SIGNING_CLIENT_SECRET}" \
				--tenant "${AZURE_SIGNING_TENANT_ID}" \
				--output none; \
			ACCESS_TOKEN=$$(az account get-access-token --resource "https://vault.azure.net" | jq -r .accessToken); \
			java -jar bin/jsign-6.0.jar \
				--storetype AZUREKEYVAULT \
				--keystore "PulumiCodeSigning" \
				--url "${AZURE_SIGNING_KEY_VAULT_URI}" \
				--storepass "$${ACCESS_TOKEN}" \
				$${file}.unsigned; \
			mv $${file}.unsigned $${file}; \
			az logout; \
		fi; \
	fi

# To make an immediately observable change to .ci-mgmt.yaml:
#
# - Edit .ci-mgmt.yaml
# - Run make ci-mgmt to apply the change locally.
#
ci-mgmt: .ci-mgmt.yaml
	go run github.com/pulumi/ci-mgmt/provider-ci@master generate
.PHONY: ci-mgmt
	fi
