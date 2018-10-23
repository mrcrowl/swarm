package source

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

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
});`

// ParseCSSFileContents parses the lines of a CSS file into bundle-ready code
func ParseCSSFileContents(name string, cssContents string, base string) (*CSSFileContents, error) {
	cssContentsWithURLsRewritten := rewriteURLStatementsInCSS(cssContents, name)
	encodedFile := jsonEncodeString(cssContentsWithURLsRewritten)
	body := fmt.Sprintf(cssTemplate, name, encodedFile)
	lines := stringToLines(body)
	return &CSSFileContents{lines}, nil
}

var rewriteURLPattern = regexp.MustCompile(`url\(['"][^'"]+['"]\)`)

func rewriteURLStatementsInCSS(css string, name string) string {
	rewrittenCSS := rewriteURLPattern.ReplaceAllStringFunc(css, func(cssURLStatement string) string {
		uri := extractURI(cssURLStatement)
		quote := string(cssURLStatement[4])
		rewrittenURI := rewriteURI(uri, name, "app")
		return "url(" + quote + rewrittenURI + quote + ")"
	})

	return rewrittenCSS
}

func rewriteURI(uri string, name string, base string) string {
	if isDataURI(uri) {
		return uri
	}

	rel, err := filepath.Rel(base, name)
	if err != nil {
		return uri
	}
	dir := path.Dir(strings.Replace(rel, "\\", "/", -1))
	rewrittenURI := path.Join(dir, uri)
	return rewrittenURI
}

// extractURI strips url(' and ') from a css url() statement
func extractURI(cssURLStatement string) string {
	return cssURLStatement[5 : len(cssURLStatement)-2]
}

func isDataURI(path string) bool {
	return strings.HasPrefix(path, "data:")
}
