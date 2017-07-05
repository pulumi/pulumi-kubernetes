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

