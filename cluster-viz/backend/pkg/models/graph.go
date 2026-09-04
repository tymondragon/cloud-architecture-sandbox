package models

// GraphNode represents a Kubernetes resource in the graph
type GraphNode struct {
	ID        string            `json:"id"`
	Kind      string            `json:"kind"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Status    string            `json:"status"`
	Labels    map[string]string `json:"labels"`
	ParentID  string            `json:"parentId,omitempty"`
	Metadata  map[string]string `json:"metadata"`
}

// GraphEdge represents a relationship between Kubernetes resources
type GraphEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Kind   string `json:"kind"` // ownership, routing, selection, dataflow, parentRef
	Label  string `json:"label,omitempty"`
}

// Graph represents the complete cluster topology
type Graph struct {
	Nodes map[string]*GraphNode `json:"nodes"`
	Edges map[string]*GraphEdge `json:"edges"`
}

// NewGraph creates an empty graph
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*GraphNode),
		Edges: make(map[string]*GraphEdge),
	}
}
