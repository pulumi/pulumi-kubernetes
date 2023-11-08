package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCrdArgs(t *testing.T) {
	parsed, err := ParseCrdArgs([]string{
		"--version",
		"1.5.3",
		"../path/to/file.yaml",
		"other/path/to/dir/",
	})
	assert.NoError(t, err)
	assert.Equal(t, "1.5.3", parsed.PackageVersion)
	assert.Equal(t, []string{
		"../path/to/file.yaml",
		"other/path/to/dir/"},
		parsed.YamlPaths,
	)
}

func TestParseCrdArgsOutOfOrder(t *testing.T) {
	parsed, err := ParseCrdArgs([]string{
		"../path/to/file.yaml",
		"--version",
		"1.5.3",
		"other/path/to/dir/",
	})
	assert.NoError(t, err)
	assert.Equal(t, "1.5.3", parsed.PackageVersion)
	assert.Equal(t, []string{
		"../path/to/file.yaml",
		"other/path/to/dir/"},
		parsed.YamlPaths,
	)
}

func TestParseCrdArgsNoVersion(t *testing.T) {
	_, err := ParseCrdArgs([]string{
		"../path/to/file.yaml",
		"other/path/to/dir/",
	})
	assert.Error(t, err)
}

func TestParseCrdArgsNoPaths(t *testing.T) {
	_, err := ParseCrdArgs([]string{
		"--version",
		"1.5.3",
	})
	assert.Error(t, err)
}
