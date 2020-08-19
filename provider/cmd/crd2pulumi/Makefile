PROJECT          := github.com/pulumi/pulumi-kubernetes/provider/v2/cmd/crd2pulumi

GO              ?= go
GOMODULE = GO111MODULE=on

ensure::
	$(GOMODULE) $(GO) mod tidy

build::
	$(GOMODULE) $(GO) build $(PROJECT)

release: rel-darwin rel-linux rel-windows

rel-darwin::
	GOOS=darwin GOARCH=386 $(GO) build -o releases/crd2pulumi-darwin-386/crd2pulumi $(PROJECT)
	tar -zcvf releases/crd2pulumi-darwin-368.tar.gz releases/crd2pulumi-darwin-386
	GOOS=darwin GOARCH=amd64 $(GO) build -o releases/crd2pulumi-darwin-amd64/crd2pulumi $(PROJECT)
	tar -zcvf releases/crd2pulumi-darwin-amd64.tar.gz releases/crd2pulumi-darwin-amd64

rel-linux::
	GOOS=linux GOARCH=386 $(GO) build -o releases/crd2pulumi-linux-386/crd2pulumi $(PROJECT)
	tar -zcvf releases/crd2pulumi-linux-386.tar.gz releases/crd2pulumi-linux-386
	GOOS=linux GOARCH=amd64 $(GO) build -o releases/crd2pulumi-linux-amd64/crd2pulumi $(PROJECT)
	tar -zcvf releases/crd2pulumi-linux-amd64.tar.gz releases/crd2pulumi-linux-amd64

rel-windows::
	GOOS=windows GOARCH=386 $(GO) build -o releases/crd2pulumi-windows-386/crd2pulumi $(PROJECT)
	tar -zcvf releases/crd2pulumi-windows-386.tar.gz releases/crd2pulumi-windows-386
	GOOS=windows GOARCH=amd64 $(GO) build -o releases/crd2pulumi-windows-amd64/crd2pulumi $(PROJECT)
	tar -zcvf releases/crd2pulumi-windows-amd64.tar.gz releases/crd2pulumi-windows-amd64
