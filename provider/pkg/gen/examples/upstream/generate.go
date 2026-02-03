package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"gopkg.in/yaml.v3"
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

	finfo, err := os.Lstat(mdPath)
	if err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(mdPath, 0600); err != nil {
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
	yaml string) string {

	return fmt.Sprintf("{{%% example %%}}\n### %s\n\n"+
		"```typescript\n%s```\n"+
		"```python\n%s```\n"+
		"```csharp\n%s```\n"+
		"```go\n%s```\n"+
		"```java\n%s```\n"+
		"```yaml\n%s```\n"+
		"{{%% /example %%}}\n",
		description, typescript, python, csharp, golang, java, yaml)
}

func processYaml(path string, mdDir string) error {
	yamlFile, err := os.Open(path)
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

		cmd := exec.Command("pulumi", "convert", "--generate-only", "--language", "typescript", "--out",
			filepath.Join(dir, "example-nodejs"))
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Dir = dir
		if err = cmd.Run(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "convert nodejs failed, ignoring: %+v", err)
		}
		content, err := os.ReadFile(filepath.Join(dir, "example-nodejs", "index.ts"))
		if err != nil {
			return err
		}
		typescript := string(content)

		cmd = exec.Command("pulumi", "convert", "--generate-only", "--language", "python", "--out",
			filepath.Join(dir, "example-py"))
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "convert python failed, ignoring: %+v", err)
		}
		content, err = os.ReadFile(filepath.Join(dir, "example-py", "__main__.py"))
		if err != nil {
			return err
		}
		python := string(content)

		cmd = exec.Command("pulumi", "convert", "--generate-only", "--language", "csharp", "--out",
			filepath.Join(dir, "example-dotnet"))
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Dir = dir
		if err = cmd.Run(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "convert dotnet failed, ignoring: %+v", err)
		}
		content, err = os.ReadFile(filepath.Join(dir, "example-dotnet", "Program.cs"))
		if err != nil {
			return err
		}
		csharp := string(content)

		cmd = exec.Command("pulumi", "convert", "--generate-only", "--language", "go", "--out",
			filepath.Join(dir, "example-go"))
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Dir = dir
		if err = cmd.Run(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "convert go failed, ignoring: %+v", err)
		}
		content, err = os.ReadFile(filepath.Join(dir, "example-go", "main.go"))
		if err != nil {
			return err
		}
		golang := string(content)

		cmd = exec.Command("pulumi", "convert", "--generate-only", "--language", "java", "--out",
			filepath.Join(dir, "example-java"))
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Dir = dir
		if err = cmd.Run(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "convert java failed, ignoring: %+v", err)
		}
		content, err = os.ReadFile(
			filepath.Join(dir, "example-java", "src", "main", "java", "generated_program", "App.java"),
		)
		if err != nil {
			return err
		}
		java := string(content)

		content, err = os.ReadFile(filepath.Join(dir, "Pulumi.yaml"))
		if err != nil {
			return err
		}
		yaml := string(content)

		exampleStrings = append(
			exampleStrings,
			markdownExample(description, typescript, python, csharp, golang, java, yaml),
		)
	}
	fmt.Fprintf(os.Stdout, "Writing %s\n", filepath.Join(mdDir, md))
	f, err := os.OpenFile(filepath.Join(mdDir, md), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer contract.IgnoreClose(f)
	_, err = f.Write([]byte(markdownExamples(exampleStrings)))
	contract.AssertNoErrorf(err, "unexpected error while writing markdown examples")
	return nil
}
