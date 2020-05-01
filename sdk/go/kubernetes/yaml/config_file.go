package yaml

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

type ConfigFile struct {
	pulumi.ResourceState

	resources pulumi.MapOutput
}

type ConfigFileArgs struct {
	// File is a path or URL that uniquely identifies a file.
	// TODO(joe): the .NET client seems to accept an input. Isn't that wrong? It will create resources inside Apply!
	File string
	// Transformations is an optional list of transformations to apply to Kubernetes resource definitions
	// before registering with the engine.
	Transformations []Transformation
	// ResourcePrefix isn optional prefix for the auto-generated resource names. For example, a resource named `bar`
	// created with resource prefix of `"foo"` would produce a resource named `"foo-bar"`.
	ResourcePrefix string
}

func NewConfigFile(ctx *pulumi.Context,
	name string, args *ConfigFileArgs, opts ...pulumi.ResourceOption) (*ConfigFile, error) {

	// Register the resulting resource state.
	configFile := &ConfigFile{
		resources: pulumi.MapOutput{},
	}
	err := ctx.RegisterComponentResource("kubernetes:yaml:ConfigFile", name, configFile, opts...)
	if err != nil {
		return nil, err
	}

	// Now provision all child resources by parsing the YAML file.
	if args != nil {
		// Make the component the parent of all subsequent resources.
		opts = append(opts, pulumi.Parent(configFile))

		// Honor the resource name prefix if specified.
		if args.ResourcePrefix != "" {
			name = args.ResourcePrefix + "-" + name
		}

		// Parse and decode the YAML files.
		rs, err := parseDecodeYamlFiles(ctx, &ConfigGroupArgs{
			Files:           []string{args.File},
			Transformations: args.Transformations,
			ResourcePrefix:  args.ResourcePrefix,
		}, true, opts...)
		if err != nil {
			return nil, err
		}
		if rs != nil {
			configFile.resources = *rs
		}

		// Finally, register all of the resources found.
		err = ctx.RegisterResourceOutputs(configFile, pulumi.Map{})
		if err != nil {
			return nil, errors.Wrapf(err, "registering child resources")
		}
	}

	return configFile, nil
}

func (cf *ConfigFile) Resources() pulumi.MapOutput {
	return cf.resources
}

func (cf *ConfigFile) GetResource(key string) pulumi.CustomResource {
	return cf.resources.MapIndex(pulumi.String(key)).(pulumi.CustomResource)
}
