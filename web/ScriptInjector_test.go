package web

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInjectInlineJavascriptModule(t *testing.T) {
	actual := InjectInlineJavascript(`<html><head><title /></head><body><div></div></body></html>`, "alert('Hello, world.');", true)
	expected := `<html><head><title /></head><body><div></div><script type="module">alert('Hello, world.');</script></body></html>`
	assert.Equal(t, expected, actual)
}

func TestInjectInlineJavascriptRegular(t *testing.T) {
	actual := InjectInlineJavascript(`<html><head><title /></head><body><div></div></body></html>`, "alert('Hello, world.');", false)
	expected := `<html><head><title /></head><body><div></div><script type="javascript">alert('Hello, world.');</script></body></html>`
	assert.Equal(t, expected, actual)
}

func TestMissingClosingBodyAppendsComment(t *testing.T) {
	html := `<html><head><title /></head><body><div></div></html>`
	actual := InjectInlineJavascript(html, "alert('Hello, world.');", true)
	assert.True(t, strings.HasPrefix(actual[len(html):], "<!--"))
	assert.True(t, strings.HasSuffix(actual, "-->"))
}
