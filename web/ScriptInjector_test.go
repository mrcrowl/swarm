package web

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInjectSrcJavascriptBodyModule(t *testing.T) {
	actual := InjectSrcJavascript(`<html><head><title /></head><body><div></div></body></html>`, "/abcd/test.js", true)
	expected := `<html><head><title /></head><body><div></div><script type="module" src="/abcd/test.js"></script></body></html>`
	assert.Equal(t, expected, actual)
}

func TestInjectSrcJavascriptHeadModule(t *testing.T) {
	actual := InjectSrcJavascript(`<html><head><title /></head><body><div></div></body></html>`, "/abcd/test.js", true)
	expected := `<html><head><title /><script type="module" src="/abcd/test.js"></script></head><body><div></div></body></html>`
	assert.Equal(t, expected, actual)
}

func TestInjectInlineJavascriptModule(t *testing.T) {
	actual := InjectInlineJavascript(`<html><head><title /></head><body><div></div></body></html>`, "alert('Hello, world.');", true)
	expected := `<html><head><title /></head><body><div></div><script type="module">alert('Hello, world.');</script></body></html>`
	assert.Equal(t, expected, actual)
}

func TestInjectInlineJavascriptRegular(t *testing.T) {
	actual := InjectInlineJavascript(`<html><head><title /></head><body><div></div></body></html>`, "alert('Hello, world.');", false)
	expected := `<html><head><title /></head><body><div></div><script type="text/javascript">alert('Hello, world.');</script></body></html>`
	assert.Equal(t, expected, actual)
}

func TestMissingClosingBodyAppendsComment(t *testing.T) {
	html := `<html><head><title /></head><body><div></div></html>`
	actual := InjectInlineJavascript(html, "alert('Hello, world.');", true)
	assert.True(t, strings.HasPrefix(actual[len(html):], "<!--"))
	assert.True(t, strings.HasSuffix(actual, "-->"))
}
