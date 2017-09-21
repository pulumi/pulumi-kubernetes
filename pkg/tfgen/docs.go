// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfgen

import (
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"strings"

	"github.com/pulumi/pulumi-fabric/pkg/diag"
	"github.com/pulumi/pulumi-fabric/pkg/util/cmdutil"
	"github.com/pulumi/pulumi-terraform/pkg/tfbridge"
)

// parsedDoc represents the data parsed from TF markdown documentation
type parsedDoc struct {
	// Description is the description of the resource
	Description string
	// Arguments includes the names and descriptions for each argument of the resource
	Arguments map[string]string
	// Attributes includes the names and descriptions for each attribute of the resource
	Attributes map[string]string
}

// getDocsForPackage extracts documentation details for the given package from TF website documentation markdown content
func getDocsForPackage(pkg string, rawname string, resinfo *tfbridge.ResourceInfo) (parsedDoc, error) {
	repo, err := getRepoDir(pkg)
	if err != nil {
		return parsedDoc{}, err
	}
	possibleMarkdownNames := []string{
		withoutPackageName(pkg, rawname) + ".html.markdown",
		withoutPackageName(pkg, rawname) + ".markdown",
		withoutPackageName(pkg, rawname) + ".html.md",
	}
	if resinfo != nil && resinfo.Docs != nil && resinfo.Docs.Source != "" {
		possibleMarkdownNames = append(possibleMarkdownNames, resinfo.Docs.Source)
	}
	markdownByts, err := readMarkdown(repo, possibleMarkdownNames)
	if err != nil {
		cmdutil.Diag().Warningf(
			diag.Message("Could not find docs for resource %v; consider overriding doc source location"), rawname)
		return parsedDoc{}, nil
	}
	doc := parseTFMarkdown(string(markdownByts), rawname)
	if resinfo != nil && resinfo.Docs != nil {
		// Merge Attributes from source into target
		mergeDocs(pkg, doc.Attributes, resinfo.Docs.IncludeAttributesFrom,
			func(s parsedDoc) map[string]string {
				return s.Attributes
			},
		)
		// Merge Arguments from source into Attributes of target
		mergeDocs(pkg, doc.Attributes, resinfo.Docs.IncludeAttributesFromArguments,
			func(s parsedDoc) map[string]string {
				return s.Arguments
			},
		)
		// Merge Arguments from source into target
		mergeDocs(pkg, doc.Arguments, resinfo.Docs.IncludeAttributesFrom,
			func(s parsedDoc) map[string]string {
				return s.Arguments
			},
		)
	}
	return doc, nil
}

// readMarkdown searches all possible locations for the markdown content
func readMarkdown(repo string, possibleLocations []string) ([]byte, error) {
	var markdownBytes []byte
	var err error
	for _, name := range possibleLocations {
		location := path.Join(repo, "website", "docs", "r", name)
		markdownBytes, err = ioutil.ReadFile(location)
		if err == nil {
			return markdownBytes, nil
		}
	}
	return nil, fmt.Errorf("Could not find markdown in any of: %v", possibleLocations)
}

// mergeDocs adds the docs specified by extractDoc from sourceFrom into the targetDocs
func mergeDocs(pkg string, targetDocs map[string]string, sourceFrom string, extractDocs func(d parsedDoc) map[string]string) error {
	if sourceFrom != "" {
		sourceDocs, err := getDocsForPackage(pkg, sourceFrom, nil)
		if err != nil {
			return err
		}
		for k, v := range extractDocs(sourceDocs) {
			targetDocs[k] = v
		}
	}
	return nil
}

var argumentBulletRegexp = regexp.MustCompile("\\*\\s+`([a-zA-z0-9_]*)`\\s+(\\([a-zA-Z]*\\)\\s*)?-?\\s+(\\([^\\)]*\\)\\s*)?(.*)")
var attributeBulletRegexp = regexp.MustCompile("\\*\\s+`([a-zA-z0-9_]*)`\\s+-\\s+(.*)")

// parseTFMarkdown takes a TF website markdown doc and extracts a structured represenation for use in
// generating doc comments
func parseTFMarkdown(markdown string, rawname string) parsedDoc {
	var ret parsedDoc
	ret.Arguments = map[string]string{}
	ret.Attributes = map[string]string{}
	sections := strings.Split(markdown, "\n## ")
	for _, section := range sections {
		lines := strings.Split(section, "\n")
		if len(lines) == 0 {
			cmdutil.Diag().Warningf(
				diag.Message("Unparseable doc section for  %v; consider overriding doc source location"), rawname)
		}
		switch lines[0] {
		case "Arguments Reference", "Argument Reference", "Nested Blocks", "Nested blocks":
			lastMatch := ""
			for _, line := range lines {
				matches := argumentBulletRegexp.FindStringSubmatch(line)
				if len(matches) >= 4 {
					// found a property bullet, extract the name and description
					ret.Arguments[matches[1]] = matches[4]
					lastMatch = matches[1]
				} else if len(matches) > 0 {
					fmt.Printf("What happened?")
				} else if strings.TrimSpace(line) != "" && lastMatch != "" {
					// this is a continuation of the previous bullet
					ret.Arguments[lastMatch] += "\n" + strings.TrimSpace(line)
				} else {
					// This is an empty line or there were no bullets yet - clear the lastMatch
					lastMatch = ""
				}
			}
		case "Attributes Reference", "Attribute Reference":
			lastMatch := ""
			for _, line := range lines {
				matches := attributeBulletRegexp.FindStringSubmatch(line)
				if len(matches) >= 2 {
					// found a property bullet, extract the name and description
					ret.Attributes[matches[1]] = matches[2]
					lastMatch = matches[1]
				} else if len(matches) > 0 {
					fmt.Printf("What happened?")
				} else if strings.TrimSpace(line) != "" && lastMatch != "" {
					// this is a continuation of the previous bullet
					ret.Attributes[lastMatch] += "\n" + strings.TrimSpace(line)
				} else {
					// This is an empty line or there were no bullets yet - clear the lastMatch
					lastMatch = ""
				}
			}
		case "---":
			// Extract the description section
			subparts := strings.Split(section, "\n# ")
			if len(subparts) != 2 {
				cmdutil.Diag().Warningf(
					diag.Message("Expected only a single H1 in markdown for resource %v"), rawname)
			}
			sublines := strings.Split(subparts[1], "\n")
			ret.Description += strings.Join(sublines[2:], "\n")
		case "Remarks":
			// Append the remarks to the description section
			ret.Description += strings.Join(lines[2:], "\n")
		default:
			// Ignore everything else - most commonly examples and imports with unpredictable section headers.
		}
	}
	return ret
}
