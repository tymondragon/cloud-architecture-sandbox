package k8s

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/tymondragon/cloud-architecture-sandbox/cluster-viz/pkg/models"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GraphBuilder maintains the in-memory graph and computes edges
type GraphBuilder struct {
	mu    sync.RWMutex
	nodes map[string]*models.GraphNode
	edges map[string]*models.GraphEdge
}

// NewGraphBuilder creates a new graph builder
func NewGraphBuilder() *GraphBuilder {
	return &GraphBuilder{
		nodes: make(map[string]*models.GraphNode),
		edges: make(map[string]*models.GraphEdge),
	}
}

// Snapshot returns a filtered copy of the current graph (lesson-relevant resources only)
func (gb *GraphBuilder) Snapshot() *models.Graph {
	gb.mu.RLock()
	defer gb.mu.RUnlock()

	// Filter nodes to only include lesson-relevant namespaces
	nodesCopy := make(map[string]*models.GraphNode)
	for k, v := range gb.nodes {
		if isLessonRelevant(v) {
			nodesCopy[k] = v
		}
	}

	// Filter edges to only include those connecting visible nodes
	edgesCopy := make(map[string]*models.GraphEdge)
	for k, v := range gb.edges {
		if _, sourceExists := nodesCopy[v.Source]; sourceExists {
			if _, targetExists := nodesCopy[v.Target]; targetExists {
				edgesCopy[k] = v
			}
		}
	}

	return &models.Graph{
		Nodes: nodesCopy,
		Edges: edgesCopy,
	}
}

// isLessonRelevant returns true if the node should be visible in the graph
func isLessonRelevant(node *models.GraphNode) bool {
	// Infrastructure namespaces to hide
	infraNamespaces := map[string]bool{
		"kube-system":          true,
		"kube-public":          true,
		"kube-node-lease":      true,
		"cert-manager":         true,
		"envoy-gateway-system": true,
		"cnpg-system":          true,
		"cluster-viz":          true,
	}

	// Hide infrastructure namespaces
	if infraNamespaces[node.Namespace] {
		return false
	}

	// Show scenario namespaces (scenario-01, scenario-02, etc.)
	if strings.HasPrefix(node.Namespace, "scenario-") {
		return true
	}

	// Show data/messaging infrastructure namespaces
	dataNamespaces := map[string]bool{
		"redpanda":  true,
		"databases": true,
		"nats":      true,
		"rabbitmq":  true,
		"redis":     true,
	}
	if dataNamespaces[node.Namespace] {
		return true
	}

	// Show envoy-gateway namespace for Gateway API resources (not the operator)
	if node.Namespace == "envoy-gateway" {
		// Only show Gateway and HTTPRoute, not internal deployments
		if node.Kind == "Gateway" || node.Kind == "HTTPRoute" {
			return true
		}
		return false
	}

	// Hide everything else
	return false
}

// AddOrUpdateNode adds or updates a node and returns the delta event
func (gb *GraphBuilder) AddOrUpdateNode(node *models.GraphNode) models.EventType {
	gb.mu.Lock()
	defer gb.mu.Unlock()

	eventType := models.EventTypeAdded
	if _, exists := gb.nodes[node.ID]; exists {
		eventType = models.EventTypeModified
	}

	gb.nodes[node.ID] = node
	return eventType
}

// DeleteNode removes a node and its associated edges
func (gb *GraphBuilder) DeleteNode(nodeID string) {
	gb.mu.Lock()
	defer gb.mu.Unlock()

	delete(gb.nodes, nodeID)

	// Remove edges connected to this node
	for edgeID, edge := range gb.edges {
		if edge.Source == nodeID || edge.Target == nodeID {
			delete(gb.edges, edgeID)
		}
	}
}

// RebuildEdges recalculates all edges based on current nodes
func (gb *GraphBuilder) RebuildEdges() []*models.GraphEdge {
	gb.mu.Lock()
	defer gb.mu.Unlock()

	// Clear existing edges
	gb.edges = make(map[string]*models.GraphEdge)

	// Rebuild all edge types
	gb.buildOwnershipEdges()
	gb.buildParentRefEdges()
	gb.buildRoutingEdges()
	gb.buildSelectionEdges()
	gb.buildDataFlowEdges()

	// Return all edges for event broadcasting
	edges := make([]*models.GraphEdge, 0, len(gb.edges))
	for _, edge := range gb.edges {
		edges = append(edges, edge)
	}
	return edges
}

// buildOwnershipEdges creates edges from ownerReferences, skipping ReplicaSets
func (gb *GraphBuilder) buildOwnershipEdges() {
	for _, node := range gb.nodes {
		ownerChain := gb.getOwnerChain(node)

		// Find the top-level owner (skip ReplicaSets)
		var topOwner string
		for _, ownerID := range ownerChain {
			if owner, exists := gb.nodes[ownerID]; exists {
				if owner.Kind != "ReplicaSet" {
					topOwner = ownerID
					break
				}
			}
		}

		// Create edge from top owner to this node
		if topOwner != "" {
			edgeID := fmt.Sprintf("ownership-%s-%s", topOwner, node.ID)
			gb.edges[edgeID] = &models.GraphEdge{
				ID:     edgeID,
				Source: topOwner,
				Target: node.ID,
				Kind:   "ownership",
				Label:  "owns",
			}
		}
	}
}

// getOwnerChain extracts owner IDs from node metadata
func (gb *GraphBuilder) getOwnerChain(node *models.GraphNode) []string {
	ownerRefsStr, ok := node.Metadata["ownerReferences"]
	if !ok || ownerRefsStr == "" {
		return nil
	}

	// Parse comma-separated owner IDs
	return strings.Split(ownerRefsStr, ",")
}

// buildParentRefEdges creates edges from HTTPRoute to Gateway via parentRefs
func (gb *GraphBuilder) buildParentRefEdges() {
	for _, node := range gb.nodes {
		if node.Kind != "HTTPRoute" {
			continue
		}

		parentRefs, ok := node.Metadata["parentRefs"]
		if !ok || parentRefs == "" {
			continue
		}

		// Parse comma-separated parent gateway IDs
		for _, parentID := range strings.Split(parentRefs, ",") {
			if _, exists := gb.nodes[parentID]; exists {
				edgeID := fmt.Sprintf("parentRef-%s-%s", parentID, node.ID)
				gb.edges[edgeID] = &models.GraphEdge{
					ID:     edgeID,
					Source: parentID,
					Target: node.ID,
					Kind:   "parentRef",
					Label:  "routes via",
				}
			}
		}
	}
}

// buildRoutingEdges creates edges from HTTPRoute to Service via backendRefs
func (gb *GraphBuilder) buildRoutingEdges() {
	for _, node := range gb.nodes {
		if node.Kind != "HTTPRoute" {
			continue
		}

		backendRefs, ok := node.Metadata["backendRefs"]
		if !ok || backendRefs == "" {
			continue
		}

		// Parse comma-separated backend service IDs
		for _, backendID := range strings.Split(backendRefs, ",") {
			if _, exists := gb.nodes[backendID]; exists {
				edgeID := fmt.Sprintf("routing-%s-%s", node.ID, backendID)
				gb.edges[edgeID] = &models.GraphEdge{
					ID:     edgeID,
					Source: node.ID,
					Target: backendID,
					Kind:   "routing",
					Label:  "routes to",
				}
			}
		}
	}
}

// buildSelectionEdges creates edges from Service to Pod via label selectors
func (gb *GraphBuilder) buildSelectionEdges() {
	for _, svcNode := range gb.nodes {
		if svcNode.Kind != "Service" {
			continue
		}

		selectorStr, ok := svcNode.Metadata["selector"]
		if !ok || selectorStr == "" {
			continue
		}

		// Parse selector as map
		selector := parseSelector(selectorStr)

		// Find matching pods
		for _, podNode := range gb.nodes {
			if podNode.Kind != "Pod" || podNode.Namespace != svcNode.Namespace {
				continue
			}

			if matchesSelector(podNode.Labels, selector) {
				edgeID := fmt.Sprintf("selection-%s-%s", svcNode.ID, podNode.ID)
				gb.edges[edgeID] = &models.GraphEdge{
					ID:     edgeID,
					Source: svcNode.ID,
					Target: podNode.ID,
					Kind:   "selection",
					Label:  "selects",
				}
			}
		}
	}
}

// buildDataFlowEdges creates edges by scanning pod environment variables for connection strings
func (gb *GraphBuilder) buildDataFlowEdges() {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`postgres://`),
		regexp.MustCompile(`nats://`),
		regexp.MustCompile(`amqp://`),
		regexp.MustCompile(`redis://`),
		regexp.MustCompile(`([a-z0-9-]+)\.svc\.cluster\.local`),
	}

	for _, podNode := range gb.nodes {
		if podNode.Kind != "Pod" {
			continue
		}

		envVars, ok := podNode.Metadata["envVars"]
		if !ok || envVars == "" {
			continue
		}

		// Scan for connection patterns
		targets := make(map[string]bool)

		for _, pattern := range patterns {
			matches := pattern.FindAllStringSubmatch(envVars, -1)
			for _, match := range matches {
				if len(match) > 1 {
					// Extract service name from .svc.cluster.local pattern
					serviceName := match[1]
					// Look for matching service in same namespace
					targetID := fmt.Sprintf("Service/%s/%s", podNode.Namespace, serviceName)
					if _, exists := gb.nodes[targetID]; exists {
						targets[targetID] = true
					}
				} else {
					// For protocol patterns, try to infer target
					protocol := strings.TrimSuffix(match[0], "://")
					gb.inferDataFlowTarget(podNode, protocol, targets)
				}
			}
		}

		// Check REDPANDA_BROKERS specifically
		if strings.Contains(envVars, "REDPANDA_BROKERS") {
			gb.inferDataFlowTarget(podNode, "redpanda", targets)
		}

		// Create edges
		for targetID := range targets {
			edgeID := fmt.Sprintf("dataflow-%s-%s", podNode.ID, targetID)
			gb.edges[edgeID] = &models.GraphEdge{
				ID:     edgeID,
				Source: podNode.ID,
				Target: targetID,
				Kind:   "dataflow",
				Label:  "connects to",
			}
		}
	}
}

// inferDataFlowTarget tries to find a service matching the protocol hint
func (gb *GraphBuilder) inferDataFlowTarget(podNode *models.GraphNode, protocol string, targets map[string]bool) {
	// Common service name patterns
	namePatterns := map[string][]string{
		"postgres":  {"postgres", "postgresql", "db"},
		"nats":      {"nats"},
		"redis":     {"redis"},
		"redpanda":  {"redpanda"},
		"rabbitmq":  {"rabbitmq"},
	}

	patterns, ok := namePatterns[protocol]
	if !ok {
		return
	}

	for _, svcNode := range gb.nodes {
		if svcNode.Kind != "Service" {
			continue
		}

		for _, pattern := range patterns {
			if strings.Contains(strings.ToLower(svcNode.Name), pattern) {
				targets[svcNode.ID] = true
				return
			}
		}
	}
}

// parseSelector converts "key1=val1,key2=val2" to map
func parseSelector(selectorStr string) map[string]string {
	selector := make(map[string]string)
	pairs := strings.Split(selectorStr, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			selector[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return selector
}

// matchesSelector checks if labels match selector
func matchesSelector(labels, selector map[string]string) bool {
	for key, value := range selector {
		if labels[key] != value {
			return false
		}
	}
	return true
}

// BuildNodeID creates a consistent node ID
func BuildNodeID(kind, namespace, name string) string {
	if namespace == "" {
		return fmt.Sprintf("%s/%s", kind, name)
	}
	return fmt.Sprintf("%s/%s/%s", kind, namespace, name)
}

// BuildOwnerID creates an owner ID from OwnerReference
func BuildOwnerID(namespace string, owner metav1.OwnerReference) string {
	return BuildNodeID(owner.Kind, namespace, owner.Name)
}

// ExtractEnvVars converts pod env vars to a searchable string
func ExtractEnvVars(containers []corev1.Container) string {
	var envStrs []string
	for _, container := range containers {
		for _, env := range container.Env {
			envStrs = append(envStrs, fmt.Sprintf("%s=%s", env.Name, env.Value))
		}
	}
	return strings.Join(envStrs, ";")
}

// HashString creates a short hash for edge IDs
func HashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash[:8])
}
