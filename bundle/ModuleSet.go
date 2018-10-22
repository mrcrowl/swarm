package bundle

import (
	"log"
	"swarm/config"
	"swarm/source"
)

// ModuleSet is
type ModuleSet struct {
	modules []*Module
}

// CreateModuleSet creates a ModuleSet from a list of NormalisedModuleDescriptions
func CreateModuleSet(moduleDescriptions []*config.NormalisedModuleDescription) *ModuleSet {
	modules := make([]*Module, len(moduleDescriptions))
	for i, descr := range moduleDescriptions {
		modules[i] = &Module{
			description:     descr,
			fileset:         source.NewEmptyFileSet(),
			entryPoints:     descr.Include,
			excludedModules: nil,
		}
	}

	set := &ModuleSet{modules}
	for _, mod := range set.modules {
		mod.attachRelatedModules(set)
	}

	set.sort()

	return set
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

func (set *ModuleSet) links() map[string][]string {
	links := make(map[string][]string)
	for _, module := range set.modules {
		links[module.Name()] = module.links()
	}
	return links
}

func (set *ModuleSet) sort() {
	sortedModules := make([]*Module, len(set.modules))
	graph := source.NewIDGraph(set.links())
	sortedNames := graph.SortTopologically(set.names())
	for i, name := range sortedNames {
		sortedModules[i] = set.getModule(name)
	}
	set.modules = sortedModules
}

// Module is
type Module struct {
	description     *config.NormalisedModuleDescription
	fileset         *source.FileSet
	entryPoints     []string
	excludedModules []*Module
}

// Name gets the name of the module
func (mod *Module) Name() string {
	return mod.description.Name
}

func (mod *Module) links() []string {
	links := make([]string, len(mod.excludedModules))
	for i, mod := range mod.excludedModules {
		links[i] = mod.Name()
	}
	return links
}

func (mod *Module) attachRelatedModules(set *ModuleSet) {
	for _, excl := range mod.description.Exclude {
		excludedModule := set.getModule(excl)
		if excludedModule == nil {
			log.Panicf("attachRelatedModules: excluded module '%s' not found", excl)
		}
		mod.excludedModules = append(mod.excludedModules, excludedModule)
	}
}
