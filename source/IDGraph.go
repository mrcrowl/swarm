package source

// IDGraph is used for sorting topologically
type IDGraph struct {
	edges        map[string][]string
	reverseEdges map[string][]string
}

func newIDGraph(links map[string][]string) *IDGraph {
	edges := make(map[string][]string)
	reverseEdges := make(map[string][]string)

	add := func(m map[string][]string, from string, to string) {
		if froms, found := m[from]; found {
			m[from] = append(froms, to)
		} else {
			m[from] = []string{to}
		}
	}

	for id, dependencyIDs := range links {
		for _, did := range dependencyIDs {
			add(edges, id, did)
			add(reverseEdges, did, id)
		}
	}

	return &IDGraph{edges, reverseEdges}
}

func (graph *IDGraph) sortTopologically(ids []string) ([]string, []string) {
	sortedIDs := make([]string, 0, len(ids))
	independentIDs := graph.identifyIndependentIDs(ids)
	dependentIDs := makeHashset(ids)
	dependentIDs.removeAll(independentIDs)

	var indieID string
	for independentIDs.nonEmpty() {
		independentIDs, indieID = independentIDs.pop()
		sortedIDs = append(sortedIDs, indieID)

		for _, dependentID := range graph.reverseEdges[indieID] {
			graph.removeDependentID(dependentID, indieID)
			if len(graph.edges[dependentID]) == 0 {
				dependentIDs.remove(dependentID)
				independentIDs = independentIDs.push(dependentID)
			}
		}
	}
	return sortedIDs, dependentIDs.ids()
}

func (graph *IDGraph) removeDependentID(dependentID string, targetID string) {
	dependencies := graph.edges[dependentID]
	for i, dependencyID := range dependencies {
		if dependencyID == targetID {
			dependencies[i] = dependencies[len(dependencies)-1]
			graph.edges[dependentID] = dependencies[:len(dependencies)-1]
			return
		}
	}
}

func (graph *IDGraph) identifyIndependentIDs(ids []string) stringStack {
	independentIDs := make([]string, 0, 256)
	for _, id := range ids {
		dependencies := graph.edges[id]
		if len(dependencies) == 0 {
			independentIDs = append(independentIDs, id)
		}
	}
	return independentIDs
}

/////////////////

type stringStack []string

func (s stringStack) nonEmpty() bool {
	return len(s) > 0
}

func (s stringStack) push(value string) stringStack {
	return append(s, value)
}

func (s stringStack) pop() (stringStack, string) {
	l := len(s)
	return s[:l-1], s[l-1]
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
