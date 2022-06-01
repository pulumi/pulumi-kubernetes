package main

import (
	"bytes"
	"fmt"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:generate go run generate.go yaml .

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <yaml source dir path> <markdown destination path>\n", os.Args[0])
		os.Exit(1)
	}
	yamlPath := os.Args[1]
	mdPath := os.Args[2]

	if !filepath.IsAbs(yamlPath) {
		cwd, err := os.Getwd()
		contract.AssertNoError(err)
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

func processYaml(path string, mdDir string) error {
	yamlFile, err := os.Open(path)
	if err != nil {
		return err
	}

	base := filepath.Base(path)
	md := strings.NewReplacer(".yaml", ".md", ".yml", ".md").Replace(base)

	buf := bytes.Buffer{}
	_, err = buf.WriteString("{{% examples %}}\n")
	_, err = buf.WriteString("## Example Usage\n")

	defer contract.IgnoreClose(yamlFile)
	decoder := yaml.NewDecoder(yamlFile)
	for {
		example := map[string]interface{}{}
		err := decoder.Decode(&example)
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		_, err = buf.WriteString("{{% example %}}\n")
		contract.AssertNoError(err)

		_, err = buf.WriteString(fmt.Sprintf("### %s\n", example["description"]))
		contract.AssertNoError(err)

		_, err = buf.WriteString("\n")
		contract.AssertNoError(err)

		err = emitExample(example, &buf)
		if err != nil {
			return err
		}
	}
	_, err = buf.WriteString("{{% /examples %}}\n")
	contract.AssertNoError(err)
	f, err := os.OpenFile(filepath.Join(mdDir, md), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer contract.IgnoreClose(f)
	_, err = f.Write(buf.Bytes())
	contract.AssertNoError(err)
	return nil
}

func emitExample(example map[string]interface{}, f io.StringWriter) error {
	dir, err := ioutil.TempDir("", "")
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
	contract.AssertNoError(src.Close())

	_, err = f.WriteString("```typescript\n")
	contract.AssertNoError(err)
	cmd := exec.Command("pulumi", "convert", "--language", "typescript", "--out",
		filepath.Join(dir, "example-nodejs"))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = dir
	if err = cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "convert nodejs failed, ignoring: %+v", err)
	}
	content, err := ioutil.ReadFile(filepath.Join(dir, "example-nodejs", "index.ts"))
	if err != nil {
		return err
	}
	_, err = f.WriteString(string(content))
	contract.AssertNoError(err)
	_, err = f.WriteString("```\n")

	_, _ = fmt.Fprint(os.Stderr, "Converting python\n")
	_, err = f.WriteString("```python\n")
	contract.AssertNoError(err)
	cmd = exec.Command("pulumi", "convert", "--language", "python", "--out",
		filepath.Join(dir, "example-py"))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "convert python failed, ignoring: %+v", err)
	}
	content, err = ioutil.ReadFile(filepath.Join(dir, "example-py", "__main__.py"))
	if err != nil {
		return err
	}
	_, err = f.WriteString(string(content))
	contract.AssertNoError(err)
	_, err = f.WriteString("```\n")

	_, err = f.WriteString("```csharp\n")
	contract.AssertNoError(err)
	cmd = exec.Command("pulumi", "convert", "--language", "csharp", "--out",
		filepath.Join(dir, "example-dotnet"))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = dir
	if err = cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "convert go failed, ignoring: %+v", err)
	}
	content, err = ioutil.ReadFile(filepath.Join(dir, "example-dotnet", "MyStack.cs"))
	if err != nil {
		return err
	}
	_, err = f.WriteString(string(content))
	contract.AssertNoError(err)
	_, err = f.WriteString("```\n")

	_, err = f.WriteString("```go\n")
	contract.AssertNoError(err)
	cmd = exec.Command("pulumi", "convert", "--language", "go", "--out",
		filepath.Join(dir, "example-go"))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = dir
	if err = cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "convert go failed, ignoring: %+v", err)
	}
	content, err = ioutil.ReadFile(filepath.Join(dir, "example-go", "main.go"))
	if err != nil {
		return err
	}
	_, err = f.WriteString(string(content))
	contract.AssertNoError(err)
	_, err = f.WriteString("```\n")

	// TODO add java when convert supports it.
	//_, err = f.WriteString("```java\n")
	//contract.AssertNoError(err)
	//cmd = exec.Command("pulumi", "convert", "--language", "java", "--out",
	//	filepath.Join(dir, "example-java"))
	//cmd.Stderr = os.Stderr
	//cmd.Stdout = os.Stdout
	//cmd.Dir = dir
	//if err = cmd.Run(); err != nil {
	//	_, _ = fmt.Fprintf(os.Stderr, "convert java failed, ignoring: %+v", err)
	//}
	//content, err = ioutil.ReadFile(filepath.Join(dir, "example-java", "Main.java"))
	//if err != nil {
	//	return err
	//}
	//_, err = f.WriteString(string(content))
	//contract.AssertNoError(err)
	//_, err = f.WriteString("```\n")

	_, err = f.WriteString("```yaml\n")
	contract.AssertNoError(err)
	content, err = ioutil.ReadFile(filepath.Join(dir, "Pulumi.yaml"))
	if err != nil {
		return err
	}
	_, err = f.WriteString(string(content))
	contract.AssertNoError(err)
	_, err = f.WriteString("```\n")
	_, err = f.WriteString("{{% /example %}}\n")
	contract.AssertNoError(err)
	return nil
}
