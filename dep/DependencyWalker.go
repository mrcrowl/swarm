package dep

import (
	"fmt"
	"strings"
	"swarm/io"
	"swarm/source"
)

// BuildFileSet creates a FileSet by following the dependency graph of an entry file
func BuildFileSet(workspace *source.Workspace, entryFileRelativePath string) *source.FileSet {
	imports, links := followDependencyChain(workspace, entryFileRelativePath)
	fileset := source.NewFileSet(imports, links, workspace)

	return fileset
}

// UpdateFileset adds dependencies for an entry file to a FileSet
func UpdateFileset(fileset *source.FileSet, modifiedFileRelativePath string) {
	// ws := fileset.Workspace()
}

func followDependencyChain(workspace *source.Workspace, entryFileRelativePath string) ([]*source.Import, []*source.DependencyLink) {
	queue := newImportQueue()
	links := make([]*source.DependencyLink, 0, 2048)

	entryFileRelativePath = strings.Replace(entryFileRelativePath, "\\", "/", -1)
	queue.pushPath(entryFileRelativePath)

	// nonrels := make(map[string]int)

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
		for _, dep := range readDependencies(file) {
			if dep.IsSolo {
				continue
			}

			depRootRelative := imp.ToRootRelativeImport(dep)
			queue.push(depRootRelative)

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

func readDependencies(file *source.File) []*source.Import {
	var line string
	var err error
	if line, err = io.ReadFirstLine(file.Filepath); err != nil {
		return nil
	}

	var filteredDeps []*source.Import
	if dependencies, ok := source.ParseRegisterDependencies(line, true); ok {
		filteredDeps = make([]*source.Import, 0, len(dependencies))
		for _, dependencyImportPath := range dependencies {
			dependencyImport := source.NewImport(dependencyImportPath)
			filteredDeps = append(filteredDeps, dependencyImport)
		}
	}

	return filteredDeps
}
