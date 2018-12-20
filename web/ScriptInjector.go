package web

import (
	"strings"
)

// InjectSrcJavascript injects a script into an HTML page, just above the closing </body> tag
func InjectSrcJavascript(html string, src string, isModule bool) string {
	injection := createSrcScript(src, isModule)
	return injectBeforeClosingBody(html, injection)
}

// InjectInlineJavascript injects a snippet of script into an HTML page, just above the closing </body> tag
func InjectInlineJavascript(html string, script string, isModule bool) string {
	injection := createInlineScript(script, isModule)
	return injectBeforeClosingBody(html, injection)
}

func injectBeforeClosingBody(html string, text string) string {
	closingBodyPos := strings.LastIndex(html, "</body>")
	if closingBodyPos > 0 {
		injectedHTML := html[:closingBodyPos] + text + html[closingBodyPos:]
		return injectedHTML
	}

	return html + "<!-- COULD NOT FIND CLOSING BODY TAG TO INJECT -->"
}

func createInlineScript(inline string, isModule bool) string {
	typ := "text/javascript"
	if isModule {
		typ = "module"
	}

	return `<script type="` + typ + `">` + inline + `</script>`
}

func createSrcScript(src string, isModule bool) string {
	typ := "text/javascript"
	if isModule {
		typ = "module"
	}

	return `<script type="` + typ + `" src="` + src + `"></script>`
}
