package source

import "fmt"

// StringFileContents describes a systemjs file
type StringFileContents struct {
	lines []string
}

// BundleLines returns a list of lines ready to include in a SystemJSBundle
func (sfc *StringFileContents) BundleLines() []string {
	return sfc.lines
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
};`

// ParseStringFileContents parses the lines of a SystemJS formatted file into the Elements
func ParseStringFileContents(name string, fileContents string) (*StringFileContents, error) {
	encodedFile := jsonEncodeString(fileContents)

	body := fmt.Sprintf(template, name, encodedFile)
	lines := stringToLines(body)
	return &StringFileContents{lines}, nil
}
