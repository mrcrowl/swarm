package source

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"swarm/util"
)

const (
	// CSSPrefix is prepended to the id attribute of each <style/> tag to ensure they are uniquely identifable
	CSSPrefix = "__swarm__css__"
)

// CSSFileContents describes a systemjs file
type CSSFileContents struct {
	lines         []string
	rawCSSContent string
}

// BundleLines returns a list of lines ready to include in a SystemJSBundle
func (cssfc *CSSFileContents) BundleLines() []string {
	return cssfc.lines
}

// RawCSSContent returns the CSS as it was originally found in the source file
func (cssfc *CSSFileContents) RawCSSContent() string {
	return cssfc.rawCSSContent
}

// SourceMappingURL returns ""
func (cssfc *CSSFileContents) SourceMappingURL() string {
	return ""
}

const cssTemplate = `System.register("%s", [], function (_export, _context) {
	"use strict";

	return {
		setters: [],
		execute: function () {
			function injectCSS(e,t){if("undefined"!=typeof document)if(n=document.querySelector("#" + CSS.escape(t)))n.childNodes[0].textContent=e;else{var n,d=document,c="appendChild";(n=d.createElement("style")).id=t,n.type="text/css",d.getElementsByTagName("head")[0][c](n),n[c](d.createTextNode(e))}}
			var css = %s;
			var id = "%s";
			injectCSS(css, id);
		}
	}
});`

// ParseCSSFileContents parses the lines of a CSS file into bundle-ready code
func ParseCSSFileContents(name string, cssContents string, base string) (*CSSFileContents, error) {
	cssContentsWithURLsRewritten := rewriteURLStatementsInCSS(cssContents, name)
	encodedFile := util.JSONEncodeString(cssContentsWithURLsRewritten)
	body := fmt.Sprintf(cssTemplate, name, encodedFile, CSSPrefix+name)
	lines := util.StringToLines(body)
	return &CSSFileContents{lines, cssContents}, nil
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
