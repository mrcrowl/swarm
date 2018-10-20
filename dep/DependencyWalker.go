package dep

import (
	"fmt"
	"gospm/systemjs"
	"path/filepath"
	"strings"
)

func followDependencyGraph(workspace *Workspace, entryFileRelativePath string) *ImportQueue {
	queue := newImportQueue()

	entryFileRelativePath = strings.Replace(entryFileRelativePath, "\\", "/", -1)
	queue.pushPath(entryFileRelativePath)

	nonrels := make(map[string]int)

	follow := func(imp *Import) {
		var file *SourceFile
		var err error

		importPath := imp.path()
		if file, err = workspace.readSourceFile(importPath); err != nil {
			println("MISSING: " + importPath)
			// println("Could not find " + rootRelativeDepPath)
			return
		}

		for _, dep := range readDependencies(file) {
			if dep.isRooted {
				if val, ok := nonrels[dep.path()]; ok {
					nonrels[dep.path()] = val + 1
				} else {
					nonrels[dep.path()] = 1
				}
			} else {
				depRootRelative := imp.toRootRelativeImport(dep)
				queue.push(depRootRelative)
			}
		}
	}

	for queue.nonEmpty() {
		if ok, relativePath := queue.pop(); ok {
			follow(relativePath)
		}
	}

	for k, v := range nonrels {
		fmt.Printf("NON-REL: %s --> %d\n", k, v)
	}

	return queue
}

func readDependencies(file *SourceFile) []*Import {
	var line string
	var err error
	if line, err = file.ReadFirstLine(); err != nil {
		return nil
	}

	dependencies := systemjs.ParseRegisterDependencies(line)
	filteredDeps := make([]*Import, 0, len(dependencies))
	for _, dependencyImportPath := range dependencies {
		ext := filepath.Ext(dependencyImportPath)
		if ext == "" {
			dependencyImport := newImport(dependencyImportPath)
			filteredDeps = append(filteredDeps, dependencyImport)
		} else {
			println("EXT-MIX: " + dependencyImportPath)
		}
	}

	return filteredDeps
}
