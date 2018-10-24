package bundle

import (
	"fmt"
	"log"
	"swarm/config"
	"swarm/dep"
	"swarm/monitor"
	"swarm/source"
)

// Module is a container for managing part of a build
type Module struct {
	description       *config.NormalisedModuleDescription
	fileset           *source.FileSet
	entryPoints       []string
	excludedModules   []*Module
	bundledJavascript string
	bundledSourcemap  string
	bundler           *Bundler
	runtimeConfig     *config.RuntimeConfig
}

// NewModule creates a new Module from a NormalisedModuleDescripion
func NewModule(ws *source.Workspace, descr *config.NormalisedModuleDescription, runtimeConfig *config.RuntimeConfig) *Module {
	entryPoints := append([]string(nil), descr.Include...)
	return &Module{
		description:       descr,
		fileset:           source.NewEmptyFileSet(ws),
		entryPoints:       entryPoints,
		excludedModules:   nil,
		bundledJavascript: "",
		bundler:           NewBundler(),
		runtimeConfig:     runtimeConfig,
	}
}

// Name gets the name of the module
func (mod *Module) Name() string {
	return mod.description.Name
}

func (mod *Module) dirty() bool {
	return mod.fileset.Dirty()
}

// PrimaryEntryPoint gets the path to the primary entry/output point
func (mod *Module) PrimaryEntryPoint() string {
	return mod.description.RelativePath
}

func (mod *Module) excludedFilesets() []*source.FileSet {
	numExcludedModules := len(mod.excludedModules)
	if numExcludedModules == 0 {
		return nil
	}

	excludedFilesets := make([]*source.FileSet, numExcludedModules)
	for i, excl := range mod.excludedModules {
		excludedFilesets[i] = excl.fileset
	}
	return excludedFilesets
}

func (mod *Module) buildInitialFileSet() {
	excludedFilesets := mod.excludedFilesets()
	fileset := dep.BuildFileSet(mod.fileset.Workspace(), mod.PrimaryEntryPoint(), excludedFilesets)
	for _, entryPoint := range mod.entryPoints {
		dep.UpdateFileset(fileset, entryPoint, excludedFilesets)
	}
	mod.fileset = fileset
}

// absorbChanges absorbs an EventChangeset, triggering artefacts to be recompiled, when necessary
func (mod *Module) absorbChanges(changes *monitor.EventChangeset) {
	excludedFilesets := mod.excludedFilesets()
	ws := mod.fileset.Workspace()
	for _, entryPoint := range changes.Changes() {
		entryPointRelativePath, ok := ws.ToRelativePath(entryPoint.AbsoluteFilepath())
		if ok {
			dep.UpdateFileset(mod.fileset, entryPointRelativePath, excludedFilesets)
		}
	}
}

func (mod *Module) generateArtefacts() {
	mod.bundledJavascript, mod.bundledSourcemap = mod.bundler.Bundle(mod.fileset, mod.runtimeConfig, mod.PrimaryEntryPoint())
	mod.fileset.ClearDirty()
	fmt.Printf("   Bundled: /%s.js (%d files)\n", mod.PrimaryEntryPoint(), mod.fileset.Count())
}

func (mod *Module) links() []string {
	links := make([]string, len(mod.excludedModules))
	for i, mod := range mod.excludedModules {
		links[i] = mod.Name()
	}
	return links
}

func (mod *Module) attachExcludedModules(set *ModuleSet) {
	for _, excl := range mod.description.Exclude {
		excludedModule := set.getModule(excl)
		if excludedModule == nil {
			log.Panicf("attachExcludedModules: excluded module '%s' not found", excl)
		}
		mod.excludedModules = append(mod.excludedModules, excludedModule)
	}
}
