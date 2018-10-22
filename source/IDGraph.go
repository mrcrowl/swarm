package source

import (
	"fmt"
	"strings"
)

// IDGraph is used for sorting topologically
type IDGraph struct {
	egressEdges  map[string][]string
	ingressEdges map[string][]string
}

// NewIDGraph creates a new IDGraph
func NewIDGraph(links map[string][]string) *IDGraph {
	egressEdges := make(map[string][]string)
	ingressEdges := make(map[string][]string)

	add := func(m map[string][]string, from string, to string) {
		if froms, found := m[from]; found {
			m[from] = append(froms, to)
		} else {
			m[from] = []string{to}
		}
	}

	for id, dependencyIDs := range links {
		for _, did := range dependencyIDs {
			add(egressEdges, id, did)
			add(ingressEdges, did, id)
		}
	}

	return &IDGraph{egressEdges, ingressEdges}
}

// SortTopologically sorts the IDs in topographical order, using the links provided to NewIDGraph
func (graph *IDGraph) SortTopologically(ids []string) []string {
	sortedIDs := make([]string, 0, len(ids))

	independentIDs := graph.identifyIndependentIDs(ids)
	dependentIDs := makeHashset(ids)
	dependentIDs.removeAll(independentIDs.stack)

	for independentIDs.nonEmpty() || dependentIDs.nonEmpty() {
		for independentIDs.nonEmpty() {
			indieID := independentIDs.pop()
			// log.Printf("Handling indie ID: %s", indieID)
			sortedIDs = append(sortedIDs, indieID)

			for _, dependentID := range graph.ingressEdges[indieID] {
				graph.removeDependentID(dependentID, indieID)
				if len(graph.egressEdges[dependentID]) == 0 {
					dependentIDs.remove(dependentID)
					independentIDs.push(dependentID)
				}
			}
		}

		if dependentIDs.nonEmpty() {
			tributeID := dependentIDs.first()
			graph.breakCycle(tributeID)

			independentIDs = graph.identifyIndependentIDs(dependentIDs.ids())
			dependentIDs.removeAll(independentIDs.stack)
		}
	}
	return sortedIDs
}

func (graph *IDGraph) removeDependentID(dependentID string, targetID string) {
	dependencies := graph.egressEdges[dependentID]
	for i, dependencyID := range dependencies {
		if dependencyID == targetID {
			dependencies[i] = dependencies[len(dependencies)-1]
			graph.egressEdges[dependentID] = dependencies[:len(dependencies)-1]
			break
		}
	}

	reverseDependencies := graph.ingressEdges[targetID]
	for j, reverseDependencyID := range reverseDependencies {
		if reverseDependencyID == dependentID {
			reverseDependencies[j] = reverseDependencies[len(reverseDependencies)-1]
			graph.ingressEdges[targetID] = reverseDependencies[:len(reverseDependencies)-1]
			break
		}
	}
}

func (graph *IDGraph) identifyIndependentIDs(ids []string) *stringStack {
	independentIDs := make([]string, 0, 256)
	for _, id := range ids {
		dependencies := graph.egressEdges[id]
		if len(dependencies) == 0 {
			independentIDs = append(independentIDs, id)
		}
	}
	return newStringStack(independentIDs)
}

func (graph *IDGraph) breakCycles(ids []string) {
	for _, id := range ids {
		graph.breakCycle(id)
	}
}

func (graph *IDGraph) breakCycle(id string) {
	// log.Printf("BEGINNING: breakCycle")
	visited := make(map[string]bool)

	var recurse func(string, int) bool
	recurse = func(idcurr string, depth int) bool {
		// indent := strings.Repeat("\t", depth)
		// log.Printf("%sVisiting %s", indent, idcurr)
		visited[idcurr] = true
		dependencyIDs := graph.egressEdges[idcurr]

		for _, depID := range dependencyIDs {
			if visited[depID] {
				// log.Printf("Breaking cycle: %s --> %s", idcurr, depID)
				graph.removeDependentID(idcurr, depID)
				return false
			}
			if !recurse(depID, depth+1) {
				return false
			}
		}
		delete(visited, idcurr)
		// log.Printf("%sLeaving %s", indent, idcurr)
		return true
	}

	recurse(id, 0)

	// log.Printf("ENDING: breakCycle")
	return
}

func (graph *IDGraph) analyseLeftoverIDs(ids []string) {
	for _, id := range ids {
		graph.analyseForCycles(id)
	}
}

func (graph *IDGraph) analyseForCycles(id string) {
	cycles := make([]cyclePath, 0, 32)
	path := make(cyclePath, 0, 32)
	var recurse func(string)
	recurse = func(idcurr string) {
		path = append(path, idcurr)
		dependencyIDs := graph.egressEdges[idcurr]

		for _, depID := range dependencyIDs {
			seenIndex := path.seenIndex(depID)
			if seenIndex >= 0 {
				path = append(path, depID)
				pathCopy := append(cyclePath(nil), path[seenIndex:]...)
				cycles = append(cycles, pathCopy)
				path = path[:len(path)-1]
				continue
			}
			recurse(depID)
		}
		path = path[:len(path)-1]
	}

	recurse(id)

	displayCycle := func(cycle cyclePath) string {
		return strings.Join(cycle, " --> ")
	}

	if len(cycles) > 0 {
		for _, c := range cycles {
			fmt.Println(displayCycle(c))
		}
	}
}

///////////////

type cyclePath []string

func (path cyclePath) seenIndex(id string) int {
	for i, seenID := range path {
		if seenID == id {
			return i
		}
	}
	return -1
}

/////////////////

type stringStack struct {
	stack []string
	index map[string]bool
}

func newStringStack(values []string) *stringStack {
	index := make(map[string]bool)
	for _, v := range values {
		index[v] = true
	}
	return &stringStack{values, index}
}

func (s *stringStack) nonEmpty() bool {
	return len(s.stack) > 0
}

func (s *stringStack) push(value string) {
	if _, found := s.index[value]; !found {
		s.stack = append(s.stack, value)
	}
}

func (s *stringStack) pop() string {
	l := len(s.stack)
	value := s.stack[l-1]
	s.stack = s.stack[:l-1]
	delete(s.index, value)
	return value
}

/////////////////

type stringHashset map[string]bool

func makeHashset(ids []string) stringHashset {
	hashset := make(map[string]bool, len(ids))
	for _, id := range ids {
		hashset[id] = true
	}
	return hashset
}

func (shs stringHashset) nonEmpty() bool {
	return len(shs) > 0
}

func (shs stringHashset) first() string {
	if shs.nonEmpty() {
		for k := range shs {
			return k
		}
	}

	return ""
}

func (shs stringHashset) removeAll(ids []string) {
	for _, id := range ids {
		delete(shs, id)
	}
}

func (shs stringHashset) remove(id string) {
	delete(shs, id)
}

func (shs stringHashset) ids() []string {
	ids := make([]string, len(shs))
	i := 0
	for id := range shs {
		ids[i] = id
		i++
	}
	return ids
}
