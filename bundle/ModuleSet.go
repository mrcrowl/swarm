package bundle

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"swarm/config"
	"swarm/dep"
	"swarm/monitor"
	"swarm/source"
	"sync"
)

// ModuleSet is
type ModuleSet struct {
	modules []*Module
	mutex   *sync.Mutex
}

// CreateModuleSet creates a ModuleSet from a list of NormalisedModuleDescriptions
func CreateModuleSet(ws *source.Workspace, moduleDescriptions []*config.NormalisedModuleDescription, runtimeConfig *config.RuntimeConfig) *ModuleSet {
	modules := make([]*Module, len(moduleDescriptions))
	for i, descr := range moduleDescriptions {
		modules[i] = NewModule(ws, descr, runtimeConfig)
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

// NotifyChanges absorbs an EventChangeset, triggering artefacts to be recompiled, when necessary
func (set *ModuleSet) NotifyChanges(changes *monitor.EventChangeset) {
	set.mutex.Lock()
	if changes != nil {
		for _, mod := range set.modules {
			mod.absorbChanges(changes)
		}
	}

	// TODO: could this be parallelised?
	for _, mod := range set.modules {
		if mod.dirty() {
			mod.generateArtefacts()
		}
	}
	set.mutex.Unlock()
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

// GenerateHTTPHandlers creates http.HandlerFunc's that will return the bundled javascript
func (set *ModuleSet) GenerateHTTPHandlers() map[string]http.HandlerFunc {
	handlers := map[string]http.HandlerFunc{}
	for _, module := range set.modules {
		// "/app/src/ep/app.js":
		handlers["/"+module.PrimaryEntryPoint()+".js"] = func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, module.bundledJavascript)
		}
	}
	return handlers
}

// Module is a container for managing part of a build
type Module struct {
	description       *config.NormalisedModuleDescription
	fileset           *source.FileSet
	entryPoints       []string
	excludedModules   []*Module
	bundledJavascript string
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
	mod.bundledJavascript = mod.bundler.Bundle(mod.fileset, mod.runtimeConfig)
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
