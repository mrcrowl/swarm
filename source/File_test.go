package source

import (
	"github.com/mrcrowl/swarm/config"
	"github.com/mrcrowl/swarm/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

var temppath string

func setup() {
	temppath = testutil.CreateTempDirWithPrefix("File_test")
}

func teardown() {
	testutil.RemoveTempDir(temppath)
}

func getSampleFile(id string, ext string, contents string) *File {
	absoluteFilepath := testutil.WriteTextFile(temppath, "blah"+ext, contents)
	return newFile(id, absoluteFilepath)
}

func TestNoImmediateLoad(t *testing.T) {
	setup()
	f := getSampleFile("abcd", ".js", "")
	assert.False(t, f.Loaded())
	teardown()
}

func TestLoadNonExistentFile(t *testing.T) {
	f := newFile("blah", "c:\\asldkfhjaksjdfh.js")
	f.EnsureLoaded(nil)
	assert.True(t, f.Loaded())
	assert.IsType(t, &FailedFileContents{}, f.contents)
	assert.Nil(t, f.BundleBody())
}

func TestEnsureLoadedPlainJS(t *testing.T) {
	setup()
	f := getSampleFile("abcd", ".js", "alert(\"hi\")")
	f.EnsureLoaded(nil)
	assert.True(t, f.Loaded())
	assert.IsType(t, &JSFileContents{}, f.contents)
	assert.Len(t, f.BundleBody(), 3)
	teardown()
}

func TestEnsureLoadedSystemJS(t *testing.T) {
	setup()
	f := getSampleFile("abcd", ".js", `System.register([], function(export, require) { 
		alert(\"hi\");
	}`)
	f.EnsureLoaded(nil)
	assert.True(t, f.Loaded())
	assert.IsType(t, &JSFileContents{}, f.contents)
	assert.Len(t, f.BundleBody(), 3)
	assert.True(t, f.contents.(*JSFileContents).isSystemJS)
	teardown()
}

func TestEnsureLoadedCSS(t *testing.T) {
	setup()
	f := getSampleFile("abcd", ".css", "body { background: green }")
	f.EnsureLoaded(config.NewRuntimeConfig("", "app"))
	assert.True(t, f.Loaded())
	assert.IsType(t, &CSSFileContents{}, f.contents)
	assert.Len(t, f.BundleBody(), 13)
	teardown()
}

func TestExt(t *testing.T) {
	setup()
	f := getSampleFile("abcd", ".css", "body { background: green }")
	assert.Equal(t, ".css", f.Ext())
	teardown()
}

var hasSourceMapCases = map[string]struct {
	id       string
	ext      string
	contents string
	expected bool
}{
	"css-with": {
		id:       "css",
		ext:      ".css",
		contents: "body { background: green }\n//# sourceMappingURL=abcd.css.map",
		expected: false, // css source mapa not yet supported in swarm
	},
	"js-with": {
		id:  "js",
		ext: ".js",
		contents: `System.register([], function (exports_1, context_1) {
var Second;
var __moduleName = context_1 && context_1.id;
return {
	setters: [],
	execute: function () {
		Second = /** @class */ (function () {
			function Second() {
			}
			Second.Go = function () {
				return "Second";
			};
			return Second;
		}());
		exports_1("Second", Second);
	}
};
});
//# sourceMappingURL=Second.js.map`,
		expected: true,
	},
	"js-without": {
		id:  "js2",
		ext: ".js",
		contents: `System.register([], function (exports_1, context_1) {
});`,
		expected: false,
	},
}

func TestSourceMap(t *testing.T) {
	for name, tc := range hasSourceMapCases {
		t.Run(name, func(t *testing.T) {
			setup()
			f := getSampleFile(tc.id, tc.ext, tc.contents)
			f.EnsureLoaded(nil)
			has := f.SourceMap(config.NewRuntimeConfig("", ""), ".")
			assert.Equal(t, tc.expected, has != nil)
			teardown()
		})
	}
}
