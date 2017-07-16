// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"os"
	"strings"
	"testing"

	"github.com/pulumi/lumi/pkg/resource"
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
	asset := resource.NewTextAsset(text)

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
}
