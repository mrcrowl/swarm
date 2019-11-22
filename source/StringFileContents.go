package source

import (
	"fmt"
	"github.com/mrcrowl/swarm/util"
)

// StringFileContents describes a systemjs file
type StringFileContents struct {
	lines []string
}

// BundleLines returns a list of lines ready to include in a SystemJSBundle
func (sfc *StringFileContents) BundleLines() []string {
	return sfc.lines
}

// SourceMappingURL returns ""
func (sfc *StringFileContents) SourceMappingURL() string {
	return ""
}

const template = `System.register("%s", [], function (_export, _context) {
	"use strict";

	var __useDefault = %s;

	return {
		setters: [],
		execute: function () {
			_export("__useDefault", __useDefault);
			_export("default", __useDefault);
		}
	}
});`

// ParseStringFileContents parses the lines of a SystemJS formatted file into the Elements
func ParseStringFileContents(name string, fileContents string) (*StringFileContents, error) {
	encodedFile := util.JSONEncodeString(fileContents)

	body := fmt.Sprintf(template, name, encodedFile)
	lines := util.StringToLines(body)
	return &StringFileContents{lines}, nil
}
