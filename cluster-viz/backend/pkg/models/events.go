package models

// EventType represents the type of graph change event
type EventType string

const (
	EventTypeSync     EventType = "SYNC"
	EventTypeAdded    EventType = "ADDED"
	EventTypeModified EventType = "MODIFIED"
	EventTypeDeleted  EventType = "DELETED"
)

// GraphEvent represents a change to the graph
type GraphEvent struct {
	Type  EventType  `json:"type"`
	Node  *GraphNode `json:"node,omitempty"`
	Edge  *GraphEdge `json:"edge,omitempty"`
	Graph *Graph     `json:"graph,omitempty"` // Only populated for SYNC events
}
