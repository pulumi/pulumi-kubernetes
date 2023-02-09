CURRENT_OS=$(shell go env GOOS)
CURRENT_OS_ARCH=$(shell echo $(CURRENT_OS)-`go env GOARCH`)

include pulumictl.version

all: pulumictl

pulumictl: cache/${CURRENT_OS_ARCH}/${PULUMICTL_VERSION}/pulumictl pulumictl.version
	# this little cp/mv dance is needed to avoid MacOS deciding the signature is incorrect and
	# refusing to run it.
	cp cache/${CURRENT_OS_ARCH}/${PULUMICTL_VERSION}/pulumictl $@-tmp
	mv $@-tmp $@
	chmod a+x $@

cache/${CURRENT_OS_ARCH}/${PULUMICTL_VERSION}/pulumictl:
	mkdir -p $(shell dirname "$@")
	curl -L "https://github.com/pulumi/pulumictl/releases/download/v$(PULUMICTL_VERSION)/pulumictl-v$(PULUMICTL_VERSION)-$(CURRENT_OS_ARCH).tar.gz" | \
		tar -xzf - -C $(shell dirname "$@") pulumictl
