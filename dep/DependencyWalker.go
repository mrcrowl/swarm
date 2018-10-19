package dep

import (
	"fmt"
	"path/filepath"
	"strings"
)

// DependencyQueue represents a queue of dependencies to process
type DependencyQueue struct {
	dependencies    []string
	dependencyIndex map[string]bool
}

func newDependencyQueue() *DependencyQueue {
	return &DependencyQueue{
		dependencies:    make([]string, 0, 2048),
		dependencyIndex: make(map[string]bool),
	}
}

// AddDependency adds a reference to a dependent file
func (ds *DependencyQueue) push(dependency string) {
	if !ds.has(dependency) {
		ds.dependencyIndex[dependency] = true
		ds.dependencies = append(ds.dependencies, dependency)
	}
}

func (ds *DependencyQueue) pop() (bool, string) {
	if ds.count() > 0 {
		dependency := ds.dependencies[0]
		ds.dependencies = ds.dependencies[1:]
		return true, dependency
	}

	return false, ""
}

// HasDependency checks for a dependent file
func (ds *DependencyQueue) has(dependency string) bool {
	if _, ok := ds.dependencyIndex[dependency]; ok {
		return true
	}
	return false
}

// NumDependencies returns the number of dependent files
func (ds *DependencyQueue) count() int {
	return len(ds.dependencies)
}

func (ds *DependencyQueue) nonEmpty() bool {
	return len(ds.dependencies) > 0
}

func followDependencyGraph(workspace *Workspace, entryFileRelativePath string) *DependencyQueue {
	queue := newDependencyQueue()

	entryFileRelativePath = strings.Replace(entryFileRelativePath, "\\", "/", -1)
	queue.push(entryFileRelativePath)

	nonrels := make(map[string]int)

	follow := func(rootRelativeDepPath string) {
		var file *SourceFile
		var err error

		if file, err = workspace.readSourceFile(rootRelativeDepPath); err != nil {
			println("MISSING: " + rootRelativeDepPath)
			// println("Could not find " + rootRelativeDepPath)
			return
		}

		for _, dep := range readDependencies(file) {
			if strings.HasPrefix(dep, "./") || strings.HasPrefix(dep, "../") {
				dependencyRootRelativePath := workspace.resolveRelativeDependency(rootRelativeDepPath, dep)
				queue.push(dependencyRootRelativePath)
			} else {
				if val, ok := nonrels[dep]; ok {
					nonrels[dep] = val + 1
				} else {
					nonrels[dep] = 1
				}
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

func readDependencies(file *SourceFile) []string {
	var line string
	var err error
	if line, err = file.ReadFirstLine(); err != nil {
		return nil
	}

	openPos := strings.Index(line, "[")
	closePos := strings.LastIndex(line, "]")
	if openPos < 0 || closePos < 0 || closePos == (openPos+1) {
		return nil
	}

	dependencySlice := line[(openPos + 1):closePos]
	dependencies := strings.Split(dependencySlice, ", ")
	filteredDeps := make([]string, 0, len(dependencies))
	for _, quotedDependency := range dependencies {
		trimmedDep := strings.Trim(quotedDependency, "\"")
		ext := filepath.Ext(trimmedDep)
		if ext == "" {
			filteredDeps = append(filteredDeps, trimmedDep)
		} else {
			println("EXT-MIX: " + trimmedDep)
		}
	}

	return filteredDeps
}
