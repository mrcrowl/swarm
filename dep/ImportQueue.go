package dep

import "gospm/source"

// importQueue represents a queue of dependencies to process
type importQueue struct {
	imports   []*source.Import
	seenIndex map[string]*source.Import
}

func newImportQueue() *importQueue {
	return &importQueue{
		imports:   make([]*source.Import, 0, 2048),
		seenIndex: make(map[string]*source.Import),
	}
}

// outputImports lists the imports that were added to the queue
func (iq *importQueue) outputImports() []*source.Import {
	outputs := make([]*source.Import, len(iq.seenIndex))
	i := 0
	for _, v := range iq.seenIndex {
		outputs[i] = v
		i++
	}
	return outputs
}

// pushPath adds a reference to an import using its path
func (iq *importQueue) pushPath(importPath string) {
	if !iq.seen(importPath) {
		iq.push(NewImport(importPath))
	}
}

// push adds a reference to an import
func (iq *importQueue) push(imp *source.Import) {
	if !iq.seen(imp.Path()) {
		iq.seenIndex[imp.Path()] = imp
		iq.imports = append(iq.imports, imp)
	}
}

func (iq *importQueue) pop() (bool, *source.Import) {
	if iq.count() > 0 {
		imp := iq.imports[0]
		iq.imports = iq.imports[1:]
		return true, imp
	}

	return false, nil
}

// seen checks whether an import has been previously seen
func (iq *importQueue) seen(imp string) bool {
	if _, ok := iq.seenIndex[imp]; ok {
		return true
	}
	return false
}

// NumDependencies returns the number of dependent files
func (iq *importQueue) count() int {
	return len(iq.imports)
}

func (iq *importQueue) nonEmpty() bool {
	return len(iq.imports) > 0
}
