package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/tymondragon/cloud-architecture-sandbox/cluster-viz/pkg/models"
)

// StubCollector is a no-op collector for development/testing
type StubCollector struct {
	collectorType string
	componentID   string
	config        models.CollectorConfig
}

// NewStubCollector creates a new stub collector
func NewStubCollector(config models.CollectorConfig) *StubCollector {
	return &StubCollector{
		collectorType: config.Type,
		componentID:   config.Component,
		config:        config,
	}
}

func (c *StubCollector) Type() string {
	return c.collectorType
}

func (c *StubCollector) Start(ctx context.Context) error {
	// No-op for now
	return nil
}

func (c *StubCollector) Stop() error {
	return nil
}

func (c *StubCollector) GetActivity(componentID string) (*ComponentActivity, error) {
	// Return empty activity
	return &ComponentActivity{
		ComponentID: componentID,
		Timestamp:   time.Now(),
		MessageRate: 0,
		Throughput:  0,
		ErrorRate:   0,
		Latency:     0,
		ActiveCount: 0,
	}, nil
}

func (c *StubCollector) TraceFlow(ctx context.Context, flowID string, traceID string) (*FlowTrace, error) {
	return nil, fmt.Errorf("flow tracing not implemented in stub collector")
}

// RedpandaCollector collects metrics from Redpanda
type RedpandaCollector struct {
	*StubCollector
	adminAPI string
	topics   []string
}

// NewRedpandaCollector creates a new Redpanda collector
func NewRedpandaCollector(config models.CollectorConfig) (Collector, error) {
	// TODO: Implement in Phase 3
	return NewStubCollector(config), nil
}

// NATSCollector collects metrics from NATS
type NATSCollector struct {
	*StubCollector
	url      string
	subjects []string
}

// NewNATSCollector creates a new NATS collector
func NewNATSCollector(config models.CollectorConfig) (Collector, error) {
	// TODO: Implement when NATS scenarios are added
	return NewStubCollector(config), nil
}

// HTTPCollector collects metrics from HTTP endpoints
type HTTPCollector struct {
	*StubCollector
	endpoint string
}

// NewHTTPCollector creates a new HTTP collector
func NewHTTPCollector(config models.CollectorConfig) (Collector, error) {
	// TODO: Implement in Phase 3
	return NewStubCollector(config), nil
}

// PostgresCollector collects metrics from PostgreSQL
type PostgresCollector struct {
	*StubCollector
	connectionString string
}

// NewPostgresCollector creates a new Postgres collector
func NewPostgresCollector(config models.CollectorConfig) (Collector, error) {
	// TODO: Implement in Phase 3
	return NewStubCollector(config), nil
}
