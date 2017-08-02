[![Build Status](https://travis-ci.com/pulumi/pulumi-terraform.svg?token=cTUUEgrxaTEGyecqJpDn&branch=master)](https://travis-ci.com/pulumi/pulumi-terraform)

# Lumi Terraform Bridge

This bridge lets Lumi leverage [Terraform](https://terraform.io)'s extensive community of resource providers for
resource create, read, update, and delete (CRUD) operations.  It is meant to give Lumi instant breadth across many
cloud providers.  Eventually we expect to customize aspects of these providers -- either with custom implementations,
and/or augmenting them with non-CRUD operations like queries, metrics, and logs -- but this approach gives us a way to
bootstrap the system very quickly, while leveraging the considerably effort that has gone into building Terraform.

## How It Works

There are two major things involved in this bridge: design-time and runtime.

At design-time, we code-generate LumiPacks by dynamic inspection of a Terraform provider's schema.  This only works for
providers that are built using static schemas.  It is possible to write Terraform providers without this, which means
the ability to create LumiPacks would not exist, but in practice all interesting providers use it.

Second, the bridge sits between Lumi's various CRUD and validation operations, and the Terraform provider's.  This
behavior also leverages the Terraform provider schema, for operations like determining which diffs will require
replacements.  Ultimately, however, all mutating runtime operations end up going through the standard dynamic plugin
interface so that we don't need to get our hands dirty with various internal and stateful representations.

## Development

### Prerequisites

Before doing any development, there are a few prerequisites to install:

* Go: https://golang.org/dl
* [Dep](https://github.com/golang/dep): `$ go get -u github.com/golang/dep/cmd/dep`
* [GoMetaLinter](https://github.com/alecthomas/gometalinter):
    - `$ go get -u github.com/alecthomas/gometalinter`
    - `$ gometalinter --install`

### Building and Testing

There is a `Makefile` in the root that builds and tests everything.

To build, ensure `$GOPATH` is set, and clone into a standard Go workspace:

    $ git clone git@github.com:pulumi/pulumi-terraform $GOPATH/src/github.com/pulumi/pulumi-terraform
    $ cd $GOPATH/src/github.com/pulumi/pulumi-terraform

Before building, you will need to ensure dependencies have been restored to your enlistment:

    $ dep ensure

At this point you can run make to build and run tests:

    $ make

This installs the `lumi-tfgen` and `lumi-tfbridge` tools into $GOPATH/bin, which may now be run provided make exited
successfully.  The `Makefile` also supports just running tests (`make test`), just running the linter (`make lint`),
just running Govet (`make vet`), and so on.  Please refer to the `Makefile` for the full list of targets.

The packages are built separately from the tools.  To generate all Lumi packages from the Terraform modules, you can
run `make gen`.  This will output the latest into the `packs/` directory, which is version controlled.  To build all of
the resulting packages, run `make packs` and, to install them, run `make install`.

### Adding a New Terraform Provider

It is relatively easy to add a new Terraform provider:

* Add a dependency on the Terraform provider:
    - `$ dep ensure github.com/terraform-providers/terraform-provider-X`
* Add a new entry to the `Providers` map in [`pkg/tfbridge/providers.go`](
  https://github.com/pulumi/pulumi-terraform/blob/master/pkg/tfbridge/providers.go).
* Add a new provider file, similar to [`pkg/tfbridge/providers_aws.go`](
  https://github.com/pulumi/pulumi-terraform/blob/master/pkg/tfbridge/providers_aws.go):
    - It statically links with the provider `github.com/terraform-providers/terraform-provider-X`;
    - There is the opportunity for optional gradual modularity, renaming, and typing, through the various maps.
* Generate the `packs/` metadata using `lumi-tfbridge`.
* Check in all of the above.

### Augmenting Auto-Generated Code w/ Overlays

The `overlays/` directory contains additional directives that the code generator obeys when creating the final
`packs/`.  Namely, any additional types, functions, or entire modules in this directory may be merged into the
resulting package.  This can be useful for helper modules and functions, in addition to gradual typing, such as using
strongly typed enums in places where Terraform may only have weakly typed strings.

To do this, first add the files in the appropriate package sub-directory, and then add the requisite directives
to the provider file.  See `overlays/aws/` for an example of this in action.

