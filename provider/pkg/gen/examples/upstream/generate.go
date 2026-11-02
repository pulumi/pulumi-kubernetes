package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"gopkg.in/yaml.v3"

	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
)

//go:generate go run generate.go yaml .

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stdout, "Usage: %s <yaml source dir path> <markdown destination path>\n", os.Args[0])
		os.Exit(1)
	}
	yamlPath := os.Args[1]
	mdPath := os.Args[2]

	if !filepath.IsAbs(yamlPath) {
		cwd, err := os.Getwd()
		contract.AssertNoErrorf(err, "unexpected error while fetching working directory")
		yamlPath = filepath.Join(cwd, yamlPath)
	}

	finfo, err := os.Lstat(mdPath) //nolint:gosec // codegen tool
	if err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(mdPath, 0600); err != nil { //nolint:gosec // codegen tool
			panic(err)
		}
	}

	if !finfo.IsDir() {
		fmt.Fprintf(os.Stderr, "Expect markdown destination %q to be a directory\n", mdPath)
		os.Exit(1)
	}

	yamls, err := os.ReadDir(yamlPath)
	if err != nil {
		panic(err)
	}
	for _, yamlFile := range yamls {
		if err := processYaml(filepath.Join(yamlPath, yamlFile.Name()), mdPath); err != nil {
			fmt.Fprintf(os.Stderr, "%+v", err)
			os.Exit(1)
		}
	}
}

func markdownExamples(examples []string) string {
	s := "{{% examples %}}\n## Example Usage\n"
	for _, example := range examples {
		s += example
	}
	s += "{{% /examples %}}\n"
	return s
}

func markdownExample(description string,
	typescript string,
	python string,
	csharp string,
	golang string,
	java string,
	yaml string,
	hcl string) string {

	return fmt.Sprintf("{{%% example %%}}\n### %s\n\n"+
		"```typescript\n%s```\n"+
		"```python\n%s```\n"+
		"```csharp\n%s```\n"+
		"```go\n%s```\n"+
		"```java\n%s```\n"+
		"```yaml\n%s```\n"+
		"```hcl\n%s```\n"+
		"{{%% /example %%}}\n",
		description, typescript, python, csharp, golang, java, yaml, hcl)
}

func processYaml(path string, mdDir string) error {
	yamlFile, err := os.Open(path) //nolint:gosec // codegen tool
	if err != nil {
		return err
	}

	base := filepath.Base(path)
	md := strings.NewReplacer(".yaml", ".md", ".yml", ".md").Replace(base)

	defer contract.IgnoreClose(yamlFile)
	decoder := yaml.NewDecoder(yamlFile)
	exampleStrings := []string{}
	for {
		example := map[string]any{}
		err := decoder.Decode(&example)
		if err == io.EOF {
			break
		}

		description := example["description"].(string)
		fmt.Fprintf(os.Stdout, "Processing %s\n", description)
		dir, err := os.MkdirTemp("", "")
		if err != nil {
			return err
		}

		defer func() {
			contract.IgnoreError(os.RemoveAll(dir))
		}()

		fmt.Fprintf(os.Stderr, "New dir: %q\n", dir)

		src, err := os.OpenFile(filepath.Join(dir, "Pulumi.yaml"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}

		if err = yaml.NewEncoder(src).Encode(example); err != nil {
			return err
		}
		contract.AssertNoErrorf(src.Close(), "unexpected error while encoding YAML")

		typescript, err := convert(dir, "typescript", "example-nodejs", "index.ts")
		if err != nil {
			return err
		}
		python, err := convert(dir, "python", "example-py", "__main__.py")
		if err != nil {
			return err
		}
		csharp, err := convert(dir, "csharp", "example-dotnet", "Program.cs")
		if err != nil {
			return err
		}
		golang, err := convert(dir, "go", "example-go", "main.go")
		if err != nil {
			return err
		}
		java, err := convert(dir, "java", "example-java",
			"src", "main", "java", "generated_program", "App.java")
		if err != nil {
			return err
		}
		hcl, err := convert(dir, "hcl", "example-hcl", "program.tf")
		if err != nil {
			return err
		}
		hcl, err = unpinProviderVersions(hcl)
		if err != nil {
			return err
		}

		content, err := os.ReadFile(filepath.Join(dir, "Pulumi.yaml"))
		if err != nil {
			return err
		}
		yaml := string(content)

		exampleStrings = append(
			exampleStrings,
			markdownExample(description, typescript, python, csharp, golang, java, yaml, hcl),
		)
	}
	fmt.Fprintf(os.Stdout, "Writing %s\n", filepath.Join(mdDir, md))
	f, err := os.OpenFile(filepath.Join(mdDir, md), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777) //nolint:gosec // codegen tool
	if err != nil {
		return err
	}
	defer contract.IgnoreClose(f)
	_, err = f.Write([]byte(markdownExamples(exampleStrings)))
	contract.AssertNoErrorf(err, "unexpected error while writing markdown examples")
	return nil
}

// convert converts the Pulumi YAML program in dir to language and returns the generated main file.
func convert(dir, language, outDir string, mainFile ...string) (string, error) {
	//nolint:gosec // we control the input
	cmd := exec.Command("pulumi", "convert", "--generate-only", "--language", language, "--out",
		filepath.Join(dir, outDir))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "convert %s failed, ignoring: %+v", language, err)
	}
	content, err := os.ReadFile(filepath.Join(append([]string{dir, outDir}, mainFile...)...))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// unpinProviderVersions removes `version` from each provider in `terraform { required_providers { ... } }`.
// The converter pins whichever provider version is installed, which would change the docs on every release.
func unpinProviderVersions(program string) (string, error) {
	f, diags := hclwrite.ParseConfig([]byte(program), "program.tf", hcl.InitialPos)
	if diags.HasErrors() {
		return "", diags
	}
	for _, tf := range f.Body().Blocks() {
		if tf.Type() != "terraform" {
			continue
		}
		for _, rp := range tf.Body().Blocks() {
			if rp.Type() != "required_providers" {
				continue
			}
			for name, attr := range rp.Body().Attributes() {
				rp.Body().SetAttributeRaw(name, dropVersionItem(attr.Expr().BuildTokens(nil)))
			}
		}
	}
	return string(hclwrite.Format(f.Bytes())), nil
}

// dropVersionItem removes the `version = ...` line from an object expression's tokens.
func dropVersionItem(tokens hclwrite.Tokens) hclwrite.Tokens {
	var out hclwrite.Tokens
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == hclsyntax.TokenIdent && string(tokens[i].Bytes) == "version" &&
			i+1 < len(tokens) && tokens[i+1].Type == hclsyntax.TokenEqual {
			for i < len(tokens) && tokens[i].Type != hclsyntax.TokenNewline {
				i++
			}
			continue
		}
		out = append(out, tokens[i])
	}
	return out
}
