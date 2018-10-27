package web

import (
	"strings"
)

// InjectInlineJavascript injects a snippet of script into an HTML page, just above the closing </body> tag
func InjectInlineJavascript(html string, script string, isModule bool) string {
	closingBodyPos := strings.LastIndex(html, "</body>")
	if closingBodyPos > 0 {
		wrappedScript := wrapInlineJavascript(script, isModule)
		modifiedHTML := html[:closingBodyPos] + wrappedScript + html[closingBodyPos:]
		return modifiedHTML
	}

	return html + "<!-- COULD NOT FIND CLOSING BODY TAG TO INJECT SCRIPT -->"
}

func wrapInlineJavascript(inline string, isModule bool) string {
	typ := "javascript"
	if isModule {
		typ = "module"
	}

	return `<script type="` + typ + `">` + inline + `</script>`
}
