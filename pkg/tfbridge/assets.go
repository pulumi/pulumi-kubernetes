// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pulumi/lumi/pkg/compiler/types/predef"
	"github.com/pulumi/lumi/pkg/resource"
	"github.com/pulumi/lumi/pkg/tokens"
	"github.com/pulumi/lumi/pkg/util/contract"
)

// AssetTranslation instructs the bridge how to translate assets into something Terraform can use.
type AssetTranslation struct {
	Kind      AssetTranslationKind   // the kind of tranlsation to perform.
	Format    resource.ArchiveFormat // an archive format, required if this is an archive.
	HashField string                 // a field to store the hash into, if any.
}

// AssetTranslationKind may be used to choose from various source and dest translation targets.
type AssetTranslationKind int

const (
	FileAsset    AssetTranslationKind = iota // turn the asset into a file on disk and pass the filename.
	BytesAsset                               // turn the asset into a []byte and pass that directly.
	FileArchive                              // turn the archive into a file on disk and pass the filename.
	BytesArchive                             // turn the asset into a []byte and pass that directly.
)

// Type fetches the Lumi runtime type corresponding to values of this asset kind.
func (a *AssetTranslation) Type() tokens.Type {
	switch a.Kind {
	case FileAsset, BytesAsset:
		return predef.LumiStdlibAssetClass
	case FileArchive, BytesArchive:
		return predef.LumiStdlibArchiveClass
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
func (a *AssetTranslation) TranslateAsset(asset resource.Asset) (interface{}, error) {
	contract.Assert(a.IsAsset())

	// TODO[pulumi/lumi#153]: support HashField.

	// Begin reading the blob.
	blob, err := asset.Read()
	if err != nil {
		return nil, err
	}
	defer contract.IgnoreClose(blob)

	// Now produce either a temp file or a binary blob, as requested.
	switch a.Kind {
	case FileAsset:
		f, err := ioutil.TempFile("", "lumi-asset-for-tf")
		if err != nil {
			return nil, err
		}
		defer f.Close()
		if _, err := io.Copy(f, blob); err != nil {
			return nil, err
		}
		return filepath.Join(os.TempDir(), f.Name()), nil
	case BytesAsset:
		return ioutil.ReadAll(blob)
	default:
		contract.Failf("Unrecognized asset translation kind: %v", a.Kind)
		return nil, nil
	}
}

// TranslateArchive translates the given archive using the directives provided by the translation info.
func (a *AssetTranslation) TranslateArchive(archive resource.Archive) (interface{}, error) {
	contract.Assert(a.IsArchive())

	// TODO[pulumi/lumi#153]: support HashField.

	// Produce either a temp file or an in-memory representation, as requested.
	switch a.Kind {
	case FileAsset:
		f, err := ioutil.TempFile("", "lumi-archive-for-tf")
		if err != nil {
			return nil, err
		}
		defer f.Close()
		if err := archive.Archive(a.Format, f); err != nil {
			return nil, err
		}
		return filepath.Join(os.TempDir(), f.Name()), nil
	case BytesAsset:
		return archive.Bytes(a.Format)
	default:
		contract.Failf("Unrecognized asset translation kind: %v", a.Kind)
		return nil, nil
	}
}
