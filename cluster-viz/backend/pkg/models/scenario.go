package models

// Scenario represents the complete architecture metadata for a learning scenario
type Scenario struct {
	Metadata          ScenarioMetadata  `json:"metadata" yaml:"metadata"`
	Patterns          ScenarioPatterns  `json:"patterns" yaml:"patterns"`
	Technologies      []Technology      `json:"technologies" yaml:"technologies"`
	Layout            LayoutConfig      `json:"layout" yaml:"layout"`
	Layers            []Layer           `json:"layers" yaml:"layers"`
	Components        []Component       `json:"components" yaml:"components"`
	Flows             []Flow            `json:"flows" yaml:"flows"`
	Telemetry         TelemetryConfig   `json:"telemetry" yaml:"telemetry"`
	DataFlows         []DataFlow        `json:"dataFlows" yaml:"dataFlows"`
	TestDataGenerator TestDataGenerator `json:"testDataGenerator" yaml:"testDataGenerator"`
}

// ScenarioMetadata contains basic information about the scenario
type ScenarioMetadata struct {
	ID          string `json:"id" yaml:"id"`
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
}

// ScenarioPatterns defines the architecture patterns demonstrated
type ScenarioPatterns struct {
	Primary   []Pattern `json:"primary" yaml:"primary"`
	Secondary []Pattern `json:"secondary" yaml:"secondary"`
}

// Pattern represents an architecture pattern
type Pattern struct {
	Name        string   `json:"name" yaml:"name"`
	Description string   `json:"description" yaml:"description"`
	Benefits    []string `json:"benefits" yaml:"benefits"`
	Color       string   `json:"color" yaml:"color"`
	Icon        string   `json:"icon" yaml:"icon"`
}

// Technology represents a technology used in the scenario
type Technology struct {
	Name        string `json:"name" yaml:"name"`
	Role        string `json:"role" yaml:"role"`
	Description string `json:"description" yaml:"description"`
}

// LayoutConfig defines how the architecture should be visualized
type LayoutConfig struct {
	Type    string        `json:"type" yaml:"type"`
	Options LayoutOptions `json:"options" yaml:"options"`
}

// LayoutOptions contains layout-specific configuration
type LayoutOptions struct {
	Spacing         string                       `json:"spacing,omitempty" yaml:"spacing,omitempty"`
	FitToViewport   bool                         `json:"fitToViewport,omitempty" yaml:"fitToViewport,omitempty"`
	CompactCards    bool                         `json:"compactCards,omitempty" yaml:"compactCards,omitempty"`
	CenterComponent string                       `json:"centerComponent,omitempty" yaml:"centerComponent,omitempty"`
	Columns         int                          `json:"columns,omitempty" yaml:"columns,omitempty"`
	Positions       map[string]map[string]float64 `json:"positions,omitempty" yaml:"positions,omitempty"`
}

// Layer represents an architectural layer
type Layer struct {
	ID          string `json:"id" yaml:"id"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Color       string `json:"color" yaml:"color"`
	Position    int    `json:"position" yaml:"position"`
}

// Component represents a logical architecture component
type Component struct {
	ID             string                 `json:"id" yaml:"id"`
	Name           string                 `json:"name" yaml:"name"`
	Type           string                 `json:"type" yaml:"type"`
	Layer          string                 `json:"layer" yaml:"layer"`
	Role           string                 `json:"role" yaml:"role"`
	Description    string                 `json:"description" yaml:"description"`
	Icon           string                 `json:"icon" yaml:"icon"`
	Implementation ComponentImplementation `json:"implementation" yaml:"implementation"`
}

// ComponentImplementation maps logical component to K8s resources
type ComponentImplementation struct {
	Kubernetes       K8sImplementation `json:"kubernetes" yaml:"kubernetes"`
	AbstractName     string            `json:"abstractName,omitempty" yaml:"abstractName,omitempty"`
	ActualTechnology string            `json:"actualTechnology,omitempty" yaml:"actualTechnology,omitempty"`
}

// K8sImplementation defines the K8s resources for a component
type K8sImplementation struct {
	Deployment  string `json:"deployment,omitempty" yaml:"deployment,omitempty"`
	StatefulSet string `json:"statefulset,omitempty" yaml:"statefulset,omitempty"`
	Service     string `json:"service,omitempty" yaml:"service,omitempty"`
	Cluster     string `json:"cluster,omitempty" yaml:"cluster,omitempty"`
	Namespace   string `json:"namespace" yaml:"namespace"`
}

// Flow represents a connection between components
type Flow struct {
	From        string    `json:"from" yaml:"from"`
	To          string    `json:"to" yaml:"to"`
	Type        string    `json:"type" yaml:"type"`
	Label       string    `json:"label" yaml:"label"`
	Description string    `json:"description" yaml:"description"`
	Pattern     string    `json:"pattern" yaml:"pattern"`
	EdgeStyle   EdgeStyle `json:"edgeStyle" yaml:"edgeStyle"`
}

// EdgeStyle defines visual appearance of an edge
type EdgeStyle struct {
	Color    string `json:"color" yaml:"color"`
	Style    string `json:"style" yaml:"style"`
	Animated bool   `json:"animated" yaml:"animated"`
}

// TelemetryConfig defines how to monitor the scenario
type TelemetryConfig struct {
	Collectors []CollectorConfig `json:"collectors" yaml:"collectors"`
}

// CollectorConfig defines a telemetry collector
type CollectorConfig struct {
	Type      string                 `json:"type" yaml:"type"`
	Component string                 `json:"component" yaml:"component"`
	Config    map[string]interface{} `json:"config" yaml:"config"`
	Metrics   []string               `json:"metrics" yaml:"metrics"`
}

// DataFlow defines an example flow for animation
type DataFlow struct {
	ID          string         `json:"id" yaml:"id"`
	Name        string         `json:"name" yaml:"name"`
	Description string         `json:"description" yaml:"description"`
	Steps       []DataFlowStep `json:"steps" yaml:"steps"`
}

// DataFlowStep represents a single step in a data flow
type DataFlowStep struct {
	Component   string                 `json:"component" yaml:"component"`
	Action      string                 `json:"action" yaml:"action"`
	Description string                 `json:"description" yaml:"description"`
	Duration    int                    `json:"duration" yaml:"duration"`
	Telemetry   *TelemetryRef          `json:"telemetry,omitempty" yaml:"telemetry,omitempty"`
	Parallel    bool                   `json:"parallel,omitempty" yaml:"parallel,omitempty"`
}

// TelemetryRef references a telemetry metric
type TelemetryRef struct {
	Collector string `json:"collector" yaml:"collector"`
	Metric    string `json:"metric" yaml:"metric"`
}

// TestDataGenerator defines how to generate test data
type TestDataGenerator struct {
	Endpoint       string           `json:"endpoint" yaml:"endpoint"`
	Method         string           `json:"method" yaml:"method"`
	SamplePayloads []SamplePayload  `json:"samplePayloads" yaml:"samplePayloads"`
}

// SamplePayload represents test data
type SamplePayload struct {
	Name         string                 `json:"name" yaml:"name"`
	Description  string                 `json:"description" yaml:"description"`
	Data         map[string]interface{} `json:"data" yaml:"data"`
	ExpectedFlow string                 `json:"expectedFlow" yaml:"expectedFlow"`
}

// LogicalGraph represents the graph built from scenario metadata
type LogicalGraph struct {
	Scenario   ScenarioMetadata         `json:"scenario"`
	Patterns   ScenarioPatterns         `json:"patterns"`
	Layers     []Layer                  `json:"layers"`
	Components map[string]GraphNode     `json:"components"`
	Edges      map[string]GraphEdge     `json:"edges"`
	Health     map[string]ComponentHealth `json:"health"`
}

// ComponentHealth represents the health of a logical component
type ComponentHealth struct {
	ComponentID string `json:"componentId"`
	Status      string `json:"status"` // "healthy", "degraded", "unhealthy", "unknown"
	Replicas    string `json:"replicas,omitempty"`
	Message     string `json:"message,omitempty"`
}
