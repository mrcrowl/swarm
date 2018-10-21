package source

import "fmt"

// CSSFileContents describes a systemjs file
type CSSFileContents struct {
	lines []string
}

// BundleLines returns a list of lines ready to include in a SystemJSBundle
func (sfc *CSSFileContents) BundleLines() []string {
	return sfc.lines
}

const cssTemplate = `System.register("%s", [], function (_export, _context) {
	"use strict";

	return {
		setters: [],
		execute: function () {
		    function injectCSS(c){if (typeof document == 'undefined') return; var d=document,a='appendChild',i='styleSheet',s=d.createElement('style');s.type='text/css';d.getElementsByTagName('head')[0][a](s);s[a](d.createTextNode(c));}
			var css = %s;
			injectCSS(css);
		}
	}
};`

// ParseCSSFileContents parses the lines of a CSS file into bundle-ready code
func ParseCSSFileContents(name string, fileContents string) (*CSSFileContents, error) {
	encodedFile := jsonEncodeString(fileContents)
	body := fmt.Sprintf(cssTemplate, name, encodedFile)
	lines := stringToLines(body)
	return &CSSFileContents{lines}, nil
}
