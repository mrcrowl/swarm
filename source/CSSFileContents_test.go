package source

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const inputCSS1 = `
.my-directive {
	background-image: url('./some-background.png');
}
`
const outputCSS1 = `
.my-directive {
	background-image: url('../common/directives/some-background.png');
}
`
const inputCSS2 = `
.my-directive {
	background-image: url("data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbm...");
}
`
const outputCSS2 = `
.my-directive {
	background-image: url("data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbm...");
}
`

func TestRewriteCSSUrls(t *testing.T) {
	cases := map[string]struct {
		css          string
		rewrittenCSS string
		base         string
	}{}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			rewrittenCSS := rewriteURLStatementsInCSS(tc.css, tc.base)
			assert.Equal(t, tc.rewrittenCSS, rewrittenCSS)
		})
	}
}

func TestRewriteCSSUrlsDataURI(t *testing.T) {
	rewrittenCSS := rewriteURLStatementsInCSS(inputCSS2, "common/directives/my-directive.css")
	assert.Equal(t, outputCSS2, rewrittenCSS)
}

func TestRewriteURI(t *testing.T) {
	cases := map[string]struct {
		uri      string
		name     string
		base     string
		expected string
	}{
		"one": {
			"fonts/ionicons.ttf?v=3.0.0",
			"app/src/ep/app.theme.css",
			"app",
			"src/ep/fonts/ionicons.ttf?v=3.0.0",
		},
		"two": {
			"../../../common/fonts/blah.png",
			"app/src/ep/app.theme.css",
			"app",
			"../common/fonts/blah.png",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actual := rewriteURI(tc.uri, tc.name, tc.base)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestIsDataURI(t *testing.T) {
	cases := map[string]bool{
		"data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbm...": true,
		"./some-background.png":                                       false,
		"../../../../../some-background.png":                          false,
	}
	for uri, expected := range cases {
		t.Run(uri, func(t *testing.T) {
			actual := isDataURI(uri)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestExtractURI(t *testing.T) {
	cases := map[string]string{
		`url('./some-background.png')`: "./some-background.png",
		`url("./some-background.png")`: "./some-background.png",
	}
	for uri, expected := range cases {
		t.Run(uri, func(t *testing.T) {
			actual := extractURI(uri)
			assert.Equal(t, expected, actual)
		})
	}
}
