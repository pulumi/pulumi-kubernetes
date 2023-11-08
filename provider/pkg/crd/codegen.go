package crd

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	pschema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
)

// GenerateFromFiles performs the entire CRD codegen process.
// The yamlPaths argument can contain both file paths and URLs.
func GenerateFromFiles(cs *CodegenSettings, yamlPaths []string) (*pschema.Package, error) {
	yamlFiles := expandFolderContent(yamlPaths)
	yamlReaders := make([]io.ReadCloser, 0, len(yamlFiles))
	for _, yamlPath := range yamlFiles {
		reader, err := ReadFromLocalOrRemote(yamlPath, map[string]string{"Accept": "application/x-yaml, text/yaml"})
		if err != nil {
			return nil, fmt.Errorf("could not open YAML document at %s: %w", yamlPath, err)
		}
		yamlReaders = append(yamlReaders, reader)
	}
	return Generate(cs, yamlReaders)
}

func expandFolderContent(yamlPaths []string) []string {
	yamlFiles := []string{}
	for _, yamlPath := range yamlPaths {
		if strings.HasPrefix(yamlPath, "https://") {
			yamlFiles = append(yamlFiles, yamlPath)
		}
		info, err := os.Stat(yamlPath)
		if err != nil {
			panic(err)
		}
		if info.IsDir() {
			entries, err := os.ReadDir(yamlPath)
			if err != nil {
				panic(err)
			}
			for _, entry := range entries {
				yamlFiles = append(yamlFiles, path.Join(yamlPath, entry.Name()))
			}
		} else {
			yamlFiles = append(yamlFiles, yamlPath)
		}
	}
	return yamlFiles
}

// Generate performs the entire CRD codegen process, reading YAML content from the given readers.
func Generate(cs *CodegenSettings, yamls []io.ReadCloser) (*pschema.Package, error) {
	// Do the actual reading of files from source, may take substantial time depending on the sources.
	pg, err := ReadPackagesFromSource(cs, yamls)
	if err != nil {
		return nil, err
	}

	// Do actual schema generation
	return pg.SchemaPackage(), nil
}
