package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadBuildDescription(t *testing.T) {
	descr, err := LoadBuildDescriptionFile("c:\\wf\\lp\\web\\App\\build\\systemjs_build_controlpanel.json")
	assert.Nil(t, err)
	assert.True(t, len(descr.Modules) > 10)
}

func TestLoadBuildDescriptionMissingSuffix(t *testing.T) {
	descr, err := LoadBuildDescriptionFile("c:\\wf\\lp\\web\\App\\build\\systemjs_build_controlpanel")
	assert.Nil(t, err)
	assert.True(t, len(descr.Modules) > 10)
}

func TestErrorLoadBuildDescriptionMissing(t *testing.T) {
	_, err := LoadBuildDescriptionFile("c:\\asdlfhasdjkf.json")
	assert.NotNil(t, err)
}

func TestErrorLoadBuildDescriptionNotJSON(t *testing.T) {
	_, err := LoadBuildDescriptionFile("c:\\wf\\lp\\web\\App\\build\\exclusion-libs-app.txt")
	assert.NotNil(t, err)
}

func TestNormaliseModules(t *testing.T) {
	build, err := LoadBuildDescriptionFile("c:\\wf\\lp\\web\\App\\build\\systemjs_build_controlpanel")
	assert.Nil(t, err)
	normMods := build.NormaliseModules("c:\\wf\\lp\\web\\App\\build")
	isNormalised := func(path string) bool {
		if strings.HasPrefix(path, "controlpanel/src") {
			return true
		}

		if !strings.HasPrefix(path, "./") && !strings.HasPrefix(path, "../") {
			return true
		}

		return false
	}
	areUnique := func(paths []string) bool {
		seen := make(map[string]bool)
		for _, path := range paths {
			if seen[path] == true {
				return false
			}
			seen[path] = true
		}
		return true
	}

	for _, normMod := range normMods {
		assert.True(t, isNormalised(normMod.RelativePath), "RelativePath is not normalised: %s", normMod.RelativePath)
		assert.True(t, areUnique(normMod.Include), "Includes not unique for: %s", normMod.Name)
		assert.True(t, areUnique(normMod.Exclude), "Excludes not unique for: %s", normMod.Name)
		for _, incl := range normMod.Include {
			assert.True(t, isNormalised(incl), "Include is not normalised: %s", incl)
		}
		for _, excl := range normMod.Include {
			assert.True(t, isNormalised(excl), "Exclude is not normalised: %s", excl)
		}
	}
}
