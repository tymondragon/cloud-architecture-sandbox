package scenario

import (
	"fmt"
	"strings"

	"github.com/tymondragon/cloud-architecture-sandbox/cluster-viz/pkg/k8s"
	"github.com/tymondragon/cloud-architecture-sandbox/cluster-viz/pkg/models"
)

// Mapper maps K8s resources to logical architecture components
type Mapper struct {
	scenario     *models.Scenario
	graphBuilder *k8s.GraphBuilder
}

// NewMapper creates a new mapper
func NewMapper(scenario *models.Scenario, graphBuilder *k8s.GraphBuilder) *Mapper {
	return &Mapper{
		scenario:     scenario,
		graphBuilder: graphBuilder,
	}
}

// BuildLogicalGraph creates a graph of logical components from K8s resources
func (m *Mapper) BuildLogicalGraph() *models.LogicalGraph {
	// Get the raw K8s graph
	k8sGraph := m.graphBuilder.Snapshot()

	// Build logical graph
	logicalGraph := &models.LogicalGraph{
		Scenario:   m.scenario.Metadata,
		Patterns:   m.scenario.Patterns,
		Layers:     m.scenario.Layers,
		Components: make(map[string]models.GraphNode),
		Edges:      make(map[string]models.GraphEdge),
		Health:     make(map[string]models.ComponentHealth),
	}

	// Map components
	for _, component := range m.scenario.Components {
		node := m.mapComponentToNode(component, k8sGraph)
		logicalGraph.Components[component.ID] = node

		// Determine component health from K8s resources
		health := m.determineComponentHealth(component, k8sGraph)
		logicalGraph.Health[component.ID] = health
	}

	// Map flows
	for _, flow := range m.scenario.Flows {
		edge := m.mapFlowToEdge(flow)
		logicalGraph.Edges[edge.ID] = edge
	}

	return logicalGraph
}

// mapComponentToNode converts a logical component to a graph node
func (m *Mapper) mapComponentToNode(component models.Component, k8sGraph *models.Graph) models.GraphNode {
	return models.GraphNode{
		ID:          component.ID,
		Kind:        component.Type,
		Name:        component.Name,
		Namespace:   component.Layer, // Use layer as "namespace" for grouping
		Status:      "Unknown",
		Labels:      map[string]string{
			"layer":       component.Layer,
			"type":        component.Type,
			"role":        component.Role,
			"icon":        component.Icon,
		},
		ParentID: component.Layer,
		Metadata: map[string]string{
			"description":      component.Description,
			"abstractName":     component.Implementation.AbstractName,
			"actualTechnology": component.Implementation.ActualTechnology,
		},
	}
}

// mapFlowToEdge converts a flow to a graph edge
func (m *Mapper) mapFlowToEdge(flow models.Flow) models.GraphEdge {
	return models.GraphEdge{
		ID:     fmt.Sprintf("%s-%s-%s", flow.Type, flow.From, flow.To),
		Source: flow.From,
		Target: flow.To,
		Kind:   flow.Type,
		Label:  flow.Label,
		Metadata: map[string]string{
			"description": flow.Description,
			"pattern":     flow.Pattern,
			"color":       flow.EdgeStyle.Color,
			"style":       flow.EdgeStyle.Style,
			"animated":    fmt.Sprintf("%t", flow.EdgeStyle.Animated),
		},
	}
}

// determineComponentHealth checks K8s resources to determine component health
func (m *Mapper) determineComponentHealth(component models.Component, k8sGraph *models.Graph) models.ComponentHealth {
	impl := component.Implementation.Kubernetes
	health := models.ComponentHealth{
		ComponentID: component.ID,
		Status:      "unknown",
	}

	// Check deployment health
	if impl.Deployment != "" {
		deploymentID := fmt.Sprintf("Deployment/%s/%s", impl.Namespace, impl.Deployment)
		if node, ok := k8sGraph.Nodes[deploymentID]; ok {
			health.Status = m.mapK8sStatus(node.Status)
			if replicas, ok := node.Metadata["replicas"]; ok {
				health.Replicas = replicas
			}
			return health
		}
	}

	// Check statefulset health
	if impl.StatefulSet != "" {
		statefulsetID := fmt.Sprintf("StatefulSet/%s/%s", impl.Namespace, impl.StatefulSet)
		if node, ok := k8sGraph.Nodes[statefulsetID]; ok {
			health.Status = m.mapK8sStatus(node.Status)
			if replicas, ok := node.Metadata["replicas"]; ok {
				health.Replicas = replicas
			}
			return health
		}
	}

	// Check cluster health (for databases)
	if impl.Cluster != "" {
		// Look for pods with cluster label
		for _, node := range k8sGraph.Nodes {
			if node.Kind == "Pod" && node.Namespace == impl.Namespace {
				if strings.Contains(node.Name, impl.Cluster) {
					health.Status = m.mapK8sStatus(node.Status)
					return health
				}
			}
		}
	}

	// Check service health
	if impl.Service != "" {
		serviceID := fmt.Sprintf("Service/%s/%s", impl.Namespace, impl.Service)
		if node, ok := k8sGraph.Nodes[serviceID]; ok {
			health.Status = m.mapK8sStatus(node.Status)
			return health
		}
	}

	return health
}

// mapK8sStatus maps K8s resource status to component health status
func (m *Mapper) mapK8sStatus(k8sStatus string) string {
	switch k8sStatus {
	case "Healthy", "Running", "Active":
		return "healthy"
	case "Degraded", "Pending":
		return "degraded"
	case "Unavailable", "Failed", "CrashLoopBackOff":
		return "unhealthy"
	default:
		return "unknown"
	}
}

// GetK8sResources returns the K8s resources for a logical component
func (m *Mapper) GetK8sResources(componentID string, k8sGraph *models.Graph) []models.GraphNode {
	// Find the component in scenario
	var component *models.Component
	for _, c := range m.scenario.Components {
		if c.ID == componentID {
			component = &c
			break
		}
	}

	if component == nil {
		return nil
	}

	impl := component.Implementation.Kubernetes
	var resources []models.GraphNode

	// Find matching K8s resources
	for _, node := range k8sGraph.Nodes {
		if node.Namespace != impl.Namespace {
			continue
		}

		matched := false
		switch {
		case impl.Deployment != "" && node.Kind == "Deployment" && node.Name == impl.Deployment:
			matched = true
		case impl.StatefulSet != "" && node.Kind == "StatefulSet" && node.Name == impl.StatefulSet:
			matched = true
		case impl.Service != "" && node.Kind == "Service" && node.Name == impl.Service:
			matched = true
		case impl.Cluster != "" && strings.Contains(node.Name, impl.Cluster):
			matched = true
		case impl.Deployment != "" && node.Kind == "Pod" && strings.Contains(node.Name, impl.Deployment):
			matched = true
		case impl.StatefulSet != "" && node.Kind == "Pod" && strings.Contains(node.Name, impl.StatefulSet):
			matched = true
		}

		if matched {
			resources = append(resources, *node)
		}
	}

	return resources
}
