MAKEFLAGS    := --warn-undefined-variables

default: pulumictl pulumi-java-gen

clean:
	rm -f pulumi-java-gen
	rm -f pulumictl

pulumi-java-gen: pulumi-java-gen.version
	@./install.sh pulumi-java-gen "$(shell cat pulumi-java-gen.version)" pulumi-java pulumi-language-java
	@touch $@

pulumictl: pulumictl.version
	@./install.sh pulumictl "$(shell cat pulumictl.version)" pulumictl pulumictl
	@touch $@
