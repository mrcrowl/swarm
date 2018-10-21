package source

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGood(t *testing.T) {
	source := `System.register(["tslib", "../../../math/katex/KatexFacade", "./FormulaInputRow"], function (exports_1, context_1) {
	var tslib_1, KatexFacade_1, FormulaInputRow_1, TouchFormulaInputRow;
	var __moduleName = context_1 && context_1.id;
	return {
		setters: [
			function (tslib_1_1) {
				tslib_1 = tslib_1_1;
			},
			function (KatexFacade_1_1) {
				KatexFacade_1 = KatexFacade_1_1;
			},
			function (FormulaInputRow_1_1) {
				FormulaInputRow_1 = FormulaInputRow_1_1;
			}
		],
		execute: function () {
		}
	}
});
//# sourceMappingURL=TouchFormulaInputRow.js.map`

	elems, err := ParseJSFileContents("abcd", source)
	assert.Nil(t, err)
	expectedImports := []string{
		`"tslib"`,
		`"../../../math/katex/KatexFacade"`,
		`"./FormulaInputRow"`,
	}
	assert.ElementsMatch(t, expectedImports, elems.imports)
	assert.Equal(t, 19, len(elems.body))
	assert.Equal(t, "TouchFormulaInputRow.js.map", elems.sourceMappingURL)
	assert.Equal(t, 20, elems.lineCount)
}

func TestParseWrongFileType(t *testing.T) {
	source := `<div>
	Hello, world.
</div>`

	elems, err := ParseJSFileContents("abcd", source)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{}, elems.imports)
	assert.Equal(t, 3, len(elems.body))
	assert.Equal(t, "", elems.sourceMappingURL)
	assert.Equal(t, 3, elems.lineCount)
	assert.False(t, elems.isSystemJS)
}

func TestParseSourceMapOnly(t *testing.T) {
	source := `<div>
	Hello, world.
</div>
//# sourceMappingURL=blah.xml`

	elems, err := ParseJSFileContents("abcd", source)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{}, elems.imports)
	assert.Equal(t, 3, len(elems.body))
	assert.Equal(t, "blah.xml", elems.sourceMappingURL)
	assert.Equal(t, 4, elems.lineCount)
	assert.False(t, elems.isSystemJS)
}

func TestParseRegisterOnly(t *testing.T) {
	source := `System.register(["tslib"], function (exports_1, context_1) {
}`

	elems, err := ParseJSFileContents("abcd", source)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(elems.imports))
	assert.Equal(t, 2, len(elems.body))
	assert.Equal(t, "", elems.sourceMappingURL)
	assert.Equal(t, 2, elems.lineCount)
}

func TestParseRegisterNoImports(t *testing.T) {
	source := `System.register([], function (exports_1, context_1) {
}`

	elems, err := ParseJSFileContents("abcd", source)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(elems.imports))
}

func TestParseRegisterInvalidRegister(t *testing.T) {
	source := `System.register(][, function (exports_1, context_1) {
}`
	elems, err := ParseJSFileContents("abcd", source)
	assert.Nil(t, err)
	assert.False(t, elems.isSystemJS)
}

func TestParseRegisterInvalidRegister2(t *testing.T) {
	source := `System.register([, function (exports_1, context_1) {
}`
	elems, err := ParseJSFileContents("abcd", source)
	assert.Nil(t, err)
	assert.False(t, elems.isSystemJS)
}

func TestParseRegisterComments1(t *testing.T) {
	source := `// some comment goes here
// another comment
System.register(["tslib"], function (exports_1, context_1) {
}`
	elems, err := ParseJSFileContents("abcd", source)
	assert.Nil(t, err)
	assert.True(t, elems.isSystemJS)
}

func TestSkipPreamblePrefix(t *testing.T) {
	source := `// comment 1
// comment 2`

	lines := strings.Split(source, "\n")
	preamble, numLines := skipPreamble(lines)
	assert.Len(t, preamble, 2)
	assert.Equal(t, 2, numLines)
}

func TestSkipPreambleBlock(t *testing.T) {
	source := `/* comment 1
 comment 2
 comment 3 */ 
blah blah something else`

	lines := strings.Split(source, "\n")
	preamble, numLines := skipPreamble(lines)
	assert.Len(t, preamble, 3)
	assert.Equal(t, 3, numLines)
}

func TestSkipPreambleMixture(t *testing.T) {
	source := `/* comment 1
 comment 2
 comment 3 */ 
// another rand comment 4
blah blah something else`

	lines := strings.Split(source, "\n")
	preamble, numLines := skipPreamble(lines)
	assert.Len(t, preamble, 4)
	assert.Equal(t, 4, numLines)
}
