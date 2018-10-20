package dep

import (
	"fmt"
	"gospm/io"
	"gospm/source"
	"path/filepath"
	"strings"
)

// BuildFileSet creates a FileSet by following the dependency graph of an entry file
func BuildFileSet(workspace *source.Workspace, entryFileRelativePath string) *source.FileSet {
	imports, links := followDependencyGraph(workspace, entryFileRelativePath)
	fileset := source.NewFileSet(workspace, imports, links)

	return fileset
}

func followDependencyGraph(workspace *source.Workspace, entryFileRelativePath string) ([]*source.Import, []*source.DependencyLink) {
	queue := newImportQueue()
	links := make([]*source.DependencyLink, 0, 2048)

	entryFileRelativePath = strings.Replace(entryFileRelativePath, "\\", "/", -1)
	queue.pushPath(entryFileRelativePath)

	nonrels := make(map[string]int)

	follow := func(imp *source.Import) {
		var file *source.File
		var err error

		importPath := imp.Path()
		if file, err = workspace.ReadSourceFile(imp); err != nil {
			println("MISSING: " + importPath)
			// println("Could not find " + rootRelativeDepPath)
			return
		}

		var dependencyIDs []string
		for _, dep := range readDependencies(file) {
			if dep.IsRooted {
				if val, ok := nonrels[dep.Path()]; ok {
					nonrels[dep.Path()] = val + 1
				} else {
					nonrels[dep.Path()] = 1
				}
			} else {
				depRootRelative := imp.ToRootRelativeImport(dep)
				queue.push(depRootRelative)

				dependencyIDs = append(dependencyIDs, dep.Path())
			}
		}

		link := source.NewDependencyLink(importPath, dependencyIDs)
		links = append(links, link)
	}

	for queue.nonEmpty() {
		if ok, relativePath := queue.pop(); ok {
			follow(relativePath)
		}
	}

	for k, v := range nonrels {
		fmt.Printf("NON-REL: %s --> %d\n", k, v)
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
	if dependencies, ok := source.ParseRegisterDependencies(line); ok {
		filteredDeps = make([]*source.Import, 0, len(dependencies))
		for _, dependencyImportPath := range dependencies {
			ext := filepath.Ext(dependencyImportPath)
			if ext == "" {
				dependencyImport := NewImport(dependencyImportPath)
				filteredDeps = append(filteredDeps, dependencyImport)
			} else {
				println("EXT-MIX: " + dependencyImportPath)
			}
		}
	}

	return filteredDeps
}
