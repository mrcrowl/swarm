package bundle

import (
	"fmt"
	"log"
	"swarm/config"
	"swarm/dep"
	"swarm/monitor"
	"swarm/source"
	"sync"
	"time"
)

// ModuleSet is
type ModuleSet struct {
	modules []*Module
	mutex   *sync.Mutex
}

// CreateModuleSet creates a ModuleSet from a list of NormalisedModuleDescriptions
func CreateModuleSet(ws *source.Workspace, moduleDescriptions []*config.NormalisedModuleDescription) *ModuleSet {
	modules := make([]*Module, len(moduleDescriptions))
	for i, descr := range moduleDescriptions {
		modules[i] = NewModule(ws, descr)
	}

	set := &ModuleSet{
		modules: modules,
		mutex:   &sync.Mutex{},
	}

	for _, mod := range set.modules {
		mod.attachExcludedModules(set)
	}

	set.sort()

	for _, mod := range set.modules {
		mod.buildInitialFileSet()
	}

	return set
}

// AbsorbChanges absorbs an EventChangeset, triggering artefacts to be recompiled, when necessary
func (set *ModuleSet) AbsorbChanges(changes *monitor.EventChangeset) {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	start := time.Now()
	for _, mod := range set.modules {
		mod.absorbChanges(changes)
	}
	defer fmt.Printf("done in %s\n", time.Since(start))
}

func (set *ModuleSet) getModule(name string) *Module {
	for _, mod := range set.modules {
		if mod.description.Name == name {
			return mod
		}
	}
	log.Panicf("getModule: could not find module with name '%s'", name)
	return nil
}

// names gets the module names (sorted topographical, assuming CreateModuleSet has finished!)
func (set *ModuleSet) names() []string {
	names := make([]string, len(set.modules))
	for i, module := range set.modules {
		names[i] = module.Name()
	}
	return names
}

func (set *ModuleSet) linksMap() map[string][]string {
	allLinks := make(map[string][]string)
	for _, module := range set.modules {
		allLinks[module.Name()] = module.links()
	}
	return allLinks
}

func (set *ModuleSet) sort() {
	sortedModules := make([]*Module, len(set.modules))
	graph := source.NewIDGraph(set.linksMap())
	sortedNames := graph.SortTopologically(set.names())
	for i, name := range sortedNames {
		sortedModules[i] = set.getModule(name)
	}
	set.modules = sortedModules
}

// Module is a container for managing part of a build
type Module struct {
	description       *config.NormalisedModuleDescription
	fileset           *source.FileSet
	entryPoints       []string
	excludedModules   []*Module
	compiledArtefacts map[string]string
	bundler           *Bundler
}

// NewModule creates a new Module from a NormalisedModuleDescripion
func NewModule(ws *source.Workspace, descr *config.NormalisedModuleDescription) *Module {
	entryPoints := []string{descr.RelativePath}
	return &Module{
		description:       descr,
		fileset:           source.NewEmptyFileSet(ws),
		entryPoints:       append(entryPoints, descr.Include...),
		excludedModules:   nil,
		compiledArtefacts: make(map[string]string),
		bundler:           NewBundler(),
	}
}

// Name gets the name of the module
func (mod *Module) Name() string {
	return mod.description.Name
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
	for _, entryPoint := range changes.Changes() {
		entryPoint := mod.fileset.Workspace().ToRelativePath(entryPoint)
		dep.UpdateFileset(mod.fileset, entryPoint)
	}

	// artefact := mod.bundler.Bundle(fileset)
	// appjs = artefact
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
