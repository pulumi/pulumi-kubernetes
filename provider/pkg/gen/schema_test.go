package gen

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/txtar"
	"sigs.k8s.io/yaml"
)

func TestPulumiSchema(t *testing.T) {
	f, err := os.ReadFile("../clients/fake/swagger.json")
	require.NoError(t, err)

	var swagger map[string]any
	err = json.Unmarshal(f, &swagger)
	require.NoError(t, err)

	_ = PulumiSchema(swagger, WithResourceOverlays(ResourceOverlays), WithTypeOverlays(TypeOverlays))
}

func TestCRDs(t *testing.T) {
	dir := filepath.Join("testdata", "crds")
	tests, err := os.ReadDir(dir)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.Name(), func(t *testing.T) {
			archive, err := txtar.ParseFile(filepath.Join(dir, tt.Name()))
			require.NoError(t, err)

			var given map[string]any
			var want string
			for _, f := range archive.Files {
				switch f.Name {
				case "given":
					err := yaml.Unmarshal(f.Data, &given)
					require.NoError(t, err, f.Name)
				case "want":
					want = string(f.Data)
				default:
					t.Fatal("unrecognized filename", f.Name)
				}
			}

			actual := PulumiSchema(given)

			actualYAML, err := yaml.Marshal(actual.Resources)
			require.NoError(t, err)

			if os.Getenv("PULUMI_ACCEPT") != "" {
				for idx, f := range archive.Files {
					if f.Name == "want" {
						archive.Files[idx].Data = actualYAML
					}
				}
				os.WriteFile(filepath.Join(dir, tt.Name()), txtar.Format(archive), 0o600)
				return
			}

			assert.YAMLEq(t, want, string(actualYAML))
		})
	}
}
