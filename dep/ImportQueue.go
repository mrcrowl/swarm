package dep

// ImportQueue represents a queue of dependencies to process
type ImportQueue struct {
	imports   []*Import
	seenIndex map[string]bool
}

func newImportQueue() *ImportQueue {
	return &ImportQueue{
		imports:   make([]*Import, 0, 2048),
		seenIndex: make(map[string]bool),
	}
}

// pushPath adds a reference to an import using its path
func (iq *ImportQueue) pushPath(importPath string) {
	if !iq.seen(importPath) {
		iq.push(newImport(importPath))
	}
}

// push adds a reference to an import
func (iq *ImportQueue) push(imp *Import) {
	if !iq.seen(imp.path()) {
		iq.seenIndex[imp.path()] = true
		iq.imports = append(iq.imports, imp)
	}
}

func (iq *ImportQueue) pop() (bool, *Import) {
	if iq.count() > 0 {
		imp := iq.imports[0]
		iq.imports = iq.imports[1:]
		return true, imp
	}

	return false, nil
}

// seen checks whether an import has been previously seen
func (iq *ImportQueue) seen(imp string) bool {
	if _, ok := iq.seenIndex[imp]; ok {
		return true
	}
	return false
}

// NumDependencies returns the number of dependent files
func (iq *ImportQueue) count() int {
	return len(iq.imports)
}

func (iq *ImportQueue) nonEmpty() bool {
	return len(iq.imports) > 0
}
