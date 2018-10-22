package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"
)

// BuildDescription describes a systemjs_build file
type BuildDescription struct {
	Modules []*ModuleDescription `json:"modules"`
	Base    string               `json:"base"`
}

// ModuleDescription describes a single module within a systemjs_build file
type ModuleDescription struct {
	Name    string   `json:"name"`
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

// NormalisedModuleDescription is a module that has paths normalised relative to the root of the workspace
type NormalisedModuleDescription struct {
	ModuleDescription
	RelativePath     string
	AbsoluteFilepath string
}

// LoadBuildDescriptionFile loads a JSON build configuration file
func LoadBuildDescriptionFile(buildFilepath string) (*BuildDescription, error) {
	if filepath.Ext(buildFilepath) == "" {
		buildFilepath += ".json"
	}

	buildBytes, e := ioutil.ReadFile(buildFilepath)
	if e != nil {
		return nil, errors.New("Invalid config file: " + buildFilepath)
	}

	jsonString := string(buildBytes)
	description, err := LoadBuildDescriptionString(jsonString)
	return description, err
}

// LoadBuildDescriptionString loads a JSON string
func LoadBuildDescriptionString(buildFileString string) (*BuildDescription, error) {
	var description *BuildDescription
	err := json.Unmarshal([]byte(buildFileString), &description)
	if err != nil {
		return nil, errors.New("Invalid JSON in config file: " + err.Error())
	}
	return description, nil
}

// NormaliseModules normalises the paths of modules relative to a root
func (build *BuildDescription) NormaliseModules(rootPath string) []*NormalisedModuleDescription {
	normalisedModules := make([]*NormalisedModuleDescription, len(build.Modules))
	for i, module := range build.Modules {
		normalisedModules[i] = module.Normalise(build.Base, rootPath)
	}
	return normalisedModules
}

// Normalise creates a NormalisedModuleDescription
func (module *ModuleDescription) Normalise(basePath string, rootPath string) *NormalisedModuleDescription {
	normaliseRelativePath := func(moduleName string) string { return path.Join(basePath, moduleName) }
	relativePath := normaliseRelativePath(module.Name)
	absoluteFilepath := filepath.Join(rootPath, basePath, module.Name)

	includes := make([]string, len(module.Include))
	for i, incl := range module.Include {
		includes[i] = normaliseRelativePath(incl)
	}

	excludes := make([]string, len(module.Exclude))
	for i, excl := range module.Exclude {
		excludes[i] = excl
	}

	return &NormalisedModuleDescription{
		ModuleDescription{
			Name:    module.Name,
			Include: includes,
			Exclude: excludes,
		},
		relativePath,
		absoluteFilepath,
	}
}
