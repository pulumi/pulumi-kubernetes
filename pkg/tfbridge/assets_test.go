// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"os"
	"strings"
	"testing"

	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	t1 := &AssetTranslation{Kind: FileAsset}
	assert.True(t, t1.IsAsset())
	assert.False(t, t1.IsArchive())
	t2 := &AssetTranslation{Kind: BytesAsset}
	assert.True(t, t2.IsAsset())
	assert.False(t, t2.IsArchive())
	t3 := &AssetTranslation{Kind: FileArchive}
	assert.False(t, t3.IsAsset())
	assert.True(t, t3.IsArchive())
	t4 := &AssetTranslation{Kind: BytesArchive}
	assert.False(t, t4.IsAsset())
	assert.True(t, t4.IsArchive())
}

func TestFileAssets(t *testing.T) {
	text := "this is a test asset"
	asset, err := resource.NewTextAsset(text)
	assert.Nil(t, err)
	asset.Hash = ""

	// First, transform the asset into a file.
	t1 := &AssetTranslation{Kind: FileAsset}
	file, err := t1.TranslateAsset(asset)
	assert.Nil(t, err)
	assert.True(t, strings.HasPrefix(file.(string), os.TempDir()))

	// Second, transform the asset into a byte blob.
	t2 := &AssetTranslation{Kind: BytesAsset}
	bytes, err := t2.TranslateAsset(asset)
	assert.Nil(t, err)
	assert.Equal(t, text, string(bytes.([]byte)))

	// Next, make sure the asset is hashed and transform it into a file.
	err = asset.EnsureHash()
	assert.Nil(t, err)
	file1, err := t1.TranslateAsset(asset)
	assert.Nil(t, err)
	assert.NotEqual(t, file.(string), file1.(string))
	assert.True(t, strings.HasSuffix(file1.(string), asset.Hash))

	// Now transform it again and ensure we get the same file.
	file2, err := t1.TranslateAsset(asset)
	assert.Nil(t, err)
	assert.Equal(t, file1, file2)

	// Now clear out the asset's contents and transform it to a file.
	asset.Text = ""
	file3, err := t1.TranslateAsset(asset)
	assert.Nil(t, err)
	assert.Equal(t, file1, file3)

	// Now clear out the asset's hash, transform the asset to a file, and ensure it has no path.
	asset.Hash = ""
	file4, err := t1.TranslateAsset(asset)
	assert.Nil(t, err)
	assert.Equal(t, file4.(string), "")
	assert.NotEqual(t, file1, file4)
}

func TestFileArchives(t *testing.T) {
	text := "this is a test asset"
	asset, err := resource.NewTextAsset(text)
	assert.Nil(t, err)

	archive, err := resource.NewAssetArchive(map[string]interface{}{"test": asset})
	assert.Nil(t, err)
	archive.Hash = ""

	// First, transform the archive into a file.
	t1 := &AssetTranslation{Kind: FileArchive, Format: resource.TarArchive}
	file, err := t1.TranslateArchive(archive)
	assert.Nil(t, err)
	assert.True(t, strings.HasPrefix(file.(string), os.TempDir()))

	// Second, transform the archive into a byte blob.
	t2 := &AssetTranslation{Kind: BytesArchive, Format: resource.TarArchive}
	_, err = t2.TranslateArchive(archive)
	assert.Nil(t, err)

	// Next, make sure the archive is hashed and transform it into a file.
	err = archive.EnsureHash()
	assert.Nil(t, err)
	file1, err := t1.TranslateArchive(archive)
	assert.Nil(t, err)
	assert.NotEqual(t, file.(string), file1.(string))
	assert.True(t, strings.HasSuffix(file1.(string), archive.Hash))

	// Now transform it again and ensure we get the same file.
	file2, err := t1.TranslateArchive(archive)
	assert.Nil(t, err)
	assert.Equal(t, file1, file2)

	// Now clear out the archive's contents and transform it to a file.
	archive.Assets = nil
	file3, err := t1.TranslateArchive(archive)
	assert.Nil(t, err)
	assert.Equal(t, file1, file3)

	// Now clear out the archive's hash, transform the archive to a file, and ensure it has no path.
	archive.Hash = ""
	file4, err := t1.TranslateArchive(archive)
	assert.Nil(t, err)
	assert.Equal(t, file4.(string), "")
	assert.NotEqual(t, file1, file4)
}
