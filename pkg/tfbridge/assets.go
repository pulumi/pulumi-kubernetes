// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"io"
	"io/ioutil"

	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/util/contract"
)

// AssetTranslation instructs the bridge how to translate assets into something Terraform can use.
type AssetTranslation struct {
	Kind      AssetTranslationKind   // the kind of translation to perform.
	Format    resource.ArchiveFormat // an archive format, required if this is an archive.
	HashField string                 // a field to store the hash into, if any.
}

// AssetTranslationKind may be used to choose from various source and dest translation targets.
type AssetTranslationKind int

const (
	// FileAsset turns the asset into a file on disk and passes the filename in its place.
	FileAsset AssetTranslationKind = iota
	// BytesAsset turns the asset into a []byte and passes it directly in-memory.
	BytesAsset
	// FileArchive turns the archive into a file on disk and passes the filename in its place.
	FileArchive
	// BytesArchive turns the asset into a []byte and passes that directly in-memory.
	BytesArchive
)

// Type fetches the Pulumi runtime type corresponding to values of this asset kind.
func (a *AssetTranslation) Type() string {
	switch a.Kind {
	case FileAsset, BytesAsset:
		return "Asset"
	case FileArchive, BytesArchive:
		return "Archive"
	default:
		contract.Failf("Unrecognized asset translation kind: %v", a.Kind)
		return ""
	}
}

// IsAsset returns true if the translation deals with an asset (rather than archive).
func (a *AssetTranslation) IsAsset() bool {
	return a.Kind == FileAsset || a.Kind == BytesAsset
}

// IsArchive returns true if the translation deals with an archive (rather than asset).
func (a *AssetTranslation) IsArchive() bool {
	return a.Kind == FileArchive || a.Kind == BytesArchive
}

// TranslateAsset translates the given asset using the directives provided by the translation info.
func (a *AssetTranslation) TranslateAsset(asset *resource.Asset) (interface{}, error) {
	contract.Assert(a.IsAsset())

	// TODO[pulumi/pulumi#153]: support HashField.

	// Begin reading the blob.
	blob, err := asset.Read()
	if err != nil {
		return nil, err
	}
	defer contract.IgnoreClose(blob)

	// Now produce either a temp file or a binary blob, as requested.
	switch a.Kind {
	case FileAsset:
		f, err := ioutil.TempFile("", "pulumi-asset")
		if err != nil {
			return nil, err
		}
		defer contract.IgnoreClose(f)
		if _, err := io.Copy(f, blob); err != nil {
			return nil, err
		}
		return f.Name(), nil
	case BytesAsset:
		return ioutil.ReadAll(blob)
	default:
		contract.Failf("Unrecognized asset translation kind: %v", a.Kind)
		return nil, nil
	}
}

// TranslateArchive translates the given archive using the directives provided by the translation info.
func (a *AssetTranslation) TranslateArchive(archive *resource.Archive) (interface{}, error) {
	contract.Assert(a.IsArchive())

	// TODO[pulumi/pulumi#153]: support HashField.

	// Produce either a temp file or an in-memory representation, as requested.
	switch a.Kind {
	case FileArchive:
		f, err := ioutil.TempFile("", "pulumi-archive")
		if err != nil {
			return nil, err
		}
		defer contract.IgnoreClose(f)
		if err := archive.Archive(a.Format, f); err != nil {
			return nil, err
		}
		return f.Name(), nil
	case BytesArchive:
		return archive.Bytes(a.Format)
	default:
		contract.Failf("Unrecognized asset translation kind: %v", a.Kind)
		return nil, nil
	}
}
