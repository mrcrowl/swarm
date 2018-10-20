package source

// DependencyLink describes a relationship between a file and its dependencies
type DependencyLink struct {
	id            string
	dependencyIDs []string
}

// NewDependencyLink creates a new DependencyLink object
func NewDependencyLink(id string, dependencyIDs []string) *DependencyLink {
	return &DependencyLink{id, dependencyIDs}
}
