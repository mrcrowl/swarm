package dep

import (
	"fmt"
	"path"
	"strings"
	"github.com/mrcrowl/swarm/source"
	"github.com/mrcrowl/swarm/util"
)

// BuildFileSet creates a FileSet by following the dependency graph of an entry file
func BuildFileSet(
	workspace *source.Workspace,
	entryFileRelativePath string,
	excludedFilesets []*source.FileSet,
	interpolationValues map[string]string,
) *source.FileSet {
	imports, links := followDependencyChain(workspace, entryFileRelativePath, excludedFilesets, interpolationValues)
	fileset := source.NewFileSet(imports, links, workspace)

	return fileset
}

// UpdateFileset adds dependencies for an entry file to a FileSet
func UpdateFileset(fileset *source.FileSet, modifiedFileRelativePath string, excludedFilesets []*source.FileSet, interpolationValues map[string]string) {
	// assume a file has been touched/changed, so:
	//
	// 1. invalidate it's content

	didRemoveJSSuffix := false
	fileID := modifiedFileRelativePath
	if path.Ext(modifiedFileRelativePath) == ".js" {
		fileID = util.RemoveExtension(modifiedFileRelativePath)
		didRemoveJSSuffix = true
	}

	file := fileset.Get(fileID)
	if file == nil && didRemoveJSSuffix {
		// maybe we're importing a .js file into a .ts file
		fileID = fileID + ".js"
		file = fileset.Get(fileID)
	}

	if file != nil {
		file.UnloadContents()
		fileset.MarkDirty()

		// 2. update the dependencies (but include "fileset" in the exclusions, so we don't follow paths we already know about)
		imports, links := followDependencyChain(fileset.Workspace(), fileID, append(excludedFilesets, fileset), interpolationValues)
		fileset.Ingest(imports, links, true)
	}
}

func followDependencyChain(
	workspace *source.Workspace,
	entryFileRelativePath string,
	excludedFilesets []*source.FileSet, /* may be nil */
	interpolationValues map[string]string,
) ([]*source.Import, []*source.DependencyLink) {
	queue := newImportQueue()
	links := make([]*source.DependencyLink, 0, 2048)

	entryFileRelativePath = strings.Replace(entryFileRelativePath, "\\", "/", -1)
	queue.pushPath(entryFileRelativePath)

	shouldEnqueue := func(dep *source.Import) bool {
		if excludedFilesets != nil {
			path := dep.Path()
			for _, exclFileset := range excludedFilesets {
				if exclFileset.Contains(path) {
					return false
				}
			}
		}
		return true
	}

	follow := func(imp *source.Import) {
		var file *source.File
		var err error

		importPath := imp.Path()
		if file, err = workspace.ReadSourceFile(imp); err != nil {
			fmt.Println("MISSING: " + importPath)
			// println("Could not find " + rootRelativeDepPath)
			return
		}

		var dependencyIDs []string
		for _, dep := range readDependencies(file, interpolationValues) {
			if dep.IsSolo {
				continue
			}

			depRootRelative := imp.ToRootRelativeImport(dep)

			if shouldEnqueue(depRootRelative) {
				queue.push(depRootRelative)
			}

			dependencyIDs = append(dependencyIDs, depRootRelative.Path())
		}

		if len(dependencyIDs) > 0 {
			link := source.NewDependencyLink(importPath, dependencyIDs)
			links = append(links, link)
		}
	}

	for queue.nonEmpty() {
		if ok, relativePath := queue.pop(); ok {
			follow(relativePath)
		}
	}

	return queue.outputImports(), links
}

func readDependencies(file *source.File, interpValues map[string]string) []*source.Import {
	var line string
	var err error
	if line, err = util.ReadFirstLine(file.Filepath); err != nil {
		return nil
	}

	var filteredDeps []*source.Import
	if dependencies, ok := source.ParseRegisterDependencies(line, true); ok {
		filteredDeps = make([]*source.Import, 0, len(dependencies))
		for _, dependencyImportPath := range dependencies {
			dependencyImport := source.NewImportWithInterpolation(dependencyImportPath, interpValues)
			filteredDeps = append(filteredDeps, dependencyImport)
		}
	}

	return filteredDeps
}
