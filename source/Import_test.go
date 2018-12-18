package source

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathConsistency(t *testing.T) {
	var sut *Import
	sut = NewImport("tslib")
	assert.Equal(t, "tslib", sut.Path())

	sut = NewImport("some/root/relative/path.js")
	assert.Equal(t, "some/root/relative/path.js", sut.Path())

	sut = NewImport("./one/abc.ts")
	assert.Equal(t, "./one/abc.ts", sut.Path())

	sut = NewImport("../one/abc.ts")
	assert.Equal(t, "../one/abc.ts", sut.Path())
}

func TestImportWithInterpolation(t *testing.T) {
	var sut = NewImportWithInterpolation("tslib-#{Hello|Hello.World}", map[string]string{"Hello.World": "Gidday"})
	assert.Equal(t, NewImport("tslib-Gidday"), sut)
}

func TestSolo(t *testing.T) {
	var sut = NewImport("tslib")
	assert.True(t, sut.IsSolo)
}

func TestRootRelativeImport(t *testing.T) {
	var sut = NewImport("app/src/ep/App")
	assert.True(t, sut.IsRooted)
	assert.False(t, sut.IsSolo)
}

func TestNotRootRelativeImportParent(t *testing.T) {
	var sut = NewImport("../../ep/App")
	assert.False(t, sut.IsRooted)
}

func TestNotRootRelativeImportSelf(t *testing.T) {
	var sut = NewImport("./App")
	assert.False(t, sut.IsRooted)
}

func TestImportExt(t *testing.T) {
	imp := NewImport("./blah/blah.mobile.less")
	assert.Equal(t, ".less", imp.Ext())
}

func TestImportContainsDirective(t *testing.T) {
	path := "import \"./login-page#{Config|Config.RELEASE_TEMPLATE_STRING}.css\";"
	assert.True(t, containsInterpolation(path))
}

func TestPerformInterpolation(t *testing.T) {
	path := "import \"./login-page#{Config|Config.RELEASE_TEMPLATE_STRING}.css\";"
	interpValues := map[string]string{
		"Config.RELEASE_TEMPLATE_STRING": ".mobile",
	}
	result := performInterpolation(path, interpValues)
	assert.Equal(t, "import \"./login-page.mobile.css\";", result)
}

func TestGetInterpolationValues(t *testing.T) {
	interpValues := readInterpolationValues("Config", []string{
		`Config.EXCEPTION_LOG_ENDPOINT = "https://logs.educationperfect.com/log";`,
		`Config.BRAND = 2 /* EDUCATION_PERFECT */;`,
		`/* ============================================== */`,
		`Config.DEBUG = true;`,
		`Config.RELATIVE_URL = "../";`,
		`Config.TEST_RELEASE = false;`,
		`Config.MOBILE_RELEASE = true;`,
		`Config.RELEASE_TEMPLATE_STRING = Config.MOBILE_RELEASE ? ".mobile" : "";`,
		`/* ============================================== */`,
	})
	assert.Contains(t, interpValues, "Config.RELEASE_TEMPLATE_STRING")
}

func TestIsJSPrimitive(t *testing.T) {
	cases := map[string]struct {
		value    string
		expected bool
	}{
		"string": {
			value:    "\"string\"",
			expected: true,
		},
		"true": {
			value:    "true",
			expected: true,
		},
		"false": {
			value:    "false",
			expected: true,
		},
		"Date": {
			value:    "new Date()",
			expected: false,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actual, _ := isJSPrimitive(tc.value)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestIsJSTernary(t *testing.T) {
	cases := map[string]struct {
		value     string
		condition string
		whenTrue  string
		whenFalse string
		expected  bool
	}{
		"existing": {
			value:     "Config.MOBILE_RELEASE ? \".mobile\" : \"\"",
			condition: "Config.MOBILE_RELEASE",
			whenTrue:  "\".mobile\"",
			whenFalse: "\"\"",
			expected:  true,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			is, condition, whenTrue, whenFalse := isJSTernary(tc.value)
			assert.Equal(t, tc.expected, is)
			if is {
				assert.Equal(t, tc.condition, condition)
				assert.Equal(t, tc.whenTrue, whenTrue)
				assert.Equal(t, tc.whenFalse, whenFalse)
			}
		})
	}
}

func TestInterpretTernary(t *testing.T) {
	cases := map[string]struct {
		condition string
		values    map[string]string
		whenTrue  string
		whenFalse string
		expected  string
	}{
		"existing": {
			values:    map[string]string{"Config.MOBILE_RELEASE": "\"true\""},
			condition: "Config.MOBILE_RELEASE",
			whenTrue:  "\".mobile\"",
			whenFalse: "\"\"",
			expected:  ".mobile",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			result := interpretTernary(tc.values, tc.condition, tc.whenTrue, tc.whenFalse)
			assert.Equal(t, tc.expected, result)
		})
	}
}
