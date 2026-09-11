package telemetry

import (
	"context"
	"time"

	"github.com/tymondragon/cloud-architecture-sandbox/cluster-viz/pkg/models"
)

// Collector collects telemetry data from a specific technology
type Collector interface {
	// Type returns the collector type (redpanda, nats, http, postgres, etc.)
	Type() string

	// Start begins collecting metrics
	Start(ctx context.Context) error

	// Stop stops collecting metrics
	Stop() error

	// GetActivity returns current activity for a component
	GetActivity(componentID string) (*ComponentActivity, error)

	// TraceFlow tracks a specific data flow by trace ID
	TraceFlow(ctx context.Context, flowID string, traceID string) (*FlowTrace, error)
}

// ComponentActivity represents real-time activity metrics for a component
type ComponentActivity struct {
	ComponentID string    `json:"componentId"`
	Timestamp   time.Time `json:"timestamp"`
	MessageRate float64   `json:"messageRate"` // messages/sec
	Throughput  float64   `json:"throughput"`  // bytes/sec
	ErrorRate   float64   `json:"errorRate"`   // errors/sec
	Latency     float64   `json:"latency"`     // milliseconds (p95)
	ActiveCount int       `json:"activeCount"` // active connections/consumers/producers
}

// FlowTrace represents a traced data flow through the system
type FlowTrace struct {
	FlowID    string     `json:"flowId"`
	TraceID   string     `json:"traceId"`
	Steps     []FlowStep `json:"steps"`
	Status    string     `json:"status"` // "active", "completed", "failed"
	StartTime time.Time  `json:"startTime"`
	Duration  int64      `json:"duration"` // milliseconds
}

// FlowStep represents a single step in a traced flow
type FlowStep struct {
	ComponentID string                 `json:"componentId"`
	Action      string                 `json:"action"`
	StartTime   time.Time              `json:"startTime"`
	Duration    int64                  `json:"duration"` // milliseconds
	Status      string                 `json:"status"`   // "active", "completed", "failed"
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Factory creates collectors based on configuration
type Factory struct {
	collectors map[string]Collector
}

// NewFactory creates a new collector factory
func NewFactory() *Factory {
	return &Factory{
		collectors: make(map[string]Collector),
	}
}

// CreateCollector creates a collector based on config
func (f *Factory) CreateCollector(config models.CollectorConfig) (Collector, error) {
	switch config.Type {
	case "redpanda":
		return NewRedpandaCollector(config)
	case "nats":
		return NewNATSCollector(config)
	case "http":
		return NewHTTPCollector(config)
	case "postgres":
		return NewPostgresCollector(config)
	default:
		// Return stub collector for unsupported types
		return NewStubCollector(config), nil
	}
}

// RegisterCollector registers a collector
func (f *Factory) RegisterCollector(componentID string, collector Collector) {
	f.collectors[componentID] = collector
}

// GetCollector retrieves a collector by component ID
func (f *Factory) GetCollector(componentID string) Collector {
	return f.collectors[componentID]
}

// StartAll starts all registered collectors
func (f *Factory) StartAll(ctx context.Context) error {
	for _, collector := range f.collectors {
		if err := collector.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}

// StopAll stops all registered collectors
func (f *Factory) StopAll() error {
	for _, collector := range f.collectors {
		if err := collector.Stop(); err != nil {
			return err
		}
	}
	return nil
}
