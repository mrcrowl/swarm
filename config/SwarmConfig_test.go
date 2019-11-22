package config

import (
	"os"
	"strings"
	"testing"

	"github.com/mrcrowl/swarm/testutil"

	"github.com/stretchr/testify/assert"
)

func TestGetDefaultRootPath(t *testing.T) {
	assert.False(t, strings.Contains(getDefaultRootPath(), "%s"))
}

const swarmConfigJSONComplete = `
{
	"root": "c:\\", 
	"monitor": {
		"extensions": [".js", ".html"],
		"debounceMillis": 350
	},
	"builds": {
		"app": {
			"path": "build/systemjs_build_app.json",
			"baseHref": "app"
		}
	},
	"server": {
		"port": 8096
	}
}
`

const swarmConfigJSONPortOnly = `
{
	"server": {
		"port": 80
	}
}
`

const swarmConfigJSONEmpty = `{}`

func TestSwarmConfig(t *testing.T) {
	cases := map[string]struct {
		json          string
		port          uint16
		debounce      uint
		numExtensions int
		numBuilds     int
		rootPath      string
	}{
		"complete": {
			json:          swarmConfigJSONComplete,
			port:          uint16(8096),
			debounce:      uint(350),
			numExtensions: 2,
			numBuilds:     1,
			rootPath:      "c:\\",
		},
		"portOnly": {
			json:          swarmConfigJSONPortOnly,
			port:          uint16(80),
			debounce:      uint(150),
			numExtensions: 4,
			numBuilds:     2,
			rootPath:      getDefaultRootPath(),
		},
		"empty": {
			json:          swarmConfigJSONEmpty,
			port:          uint16(8096),
			debounce:      uint(150),
			numExtensions: 4,
			numBuilds:     2,
			rootPath:      getDefaultRootPath(),
		},
	}

	cwd, _ := os.Getwd()
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			value, err := LoadSwarmConfigString(tc.json, cwd)
			assert.Nil(t, err)

			assert.Equal(t, tc.rootPath, value.RootPath)

			assert.Len(t, value.Monitor.Extensions, tc.numExtensions)
			assert.Equal(t, tc.debounce, value.Monitor.DebounceMillis)

			assert.Len(t, value.Builds, tc.numBuilds)

			assert.Equal(t, tc.port, value.Server.Port)
		})
	}
}

func TestNormalisePaths(t *testing.T) {
	config := &SwarmConfig{
		RootPath: "c:\\random\\",
		Builds: map[string]*RuntimeConfig{
			"one": NewRuntimeConfig("..\\building", ""),
		},
	}
	config.expandAndNormalisePaths("c:\\")
	assert.Equal(t, "c:\\random", config.RootPath)
	assert.Equal(t, "c:\\building", config.Builds["one"].BuildPath)
}

func TestTryLoadSwarmConfigFromCWD(t *testing.T) {
	temppath := testutil.CreateTempDir()
	defer testutil.RemoveTempDir(temppath)
	testutil.WriteTextFile(temppath, swarmConfigDefaultFilename, swarmConfigJSONComplete)
	os.Chdir(temppath)
	_, err := TryLoadSwarmConfigFromCWD(nil)
	assert.Nil(t, err)
	var port = uint16(1234)
	conf, _ := TryLoadSwarmConfigFromCWD(&port)
	assert.Equal(t, uint16(1234), conf.Server.Port)
}
