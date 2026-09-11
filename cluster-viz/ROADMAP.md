# Cluster-Viz: Development Roadmap

This document outlines the development milestones for the Cluster-Viz architecture pattern visualizer.

## Current Status: Milestone 2 (Interactive Features) ✅ Partially Complete

---

## Milestone 1: Core Visualization ✅ COMPLETE

**Goal**: Build foundation for K8s resource visualization with real-time updates

**Completed Features**:
- ✅ Backend K8s client and resource watching (9+ resource types)
- ✅ Physical graph builder with edge inference
- ✅ WebSocket server for real-time updates
- ✅ REST API for graph snapshots
- ✅ Frontend Svelte 4 application
- ✅ Custom SVG-based rendering (no external libraries)
- ✅ Dark theme visual design
- ✅ Multi-stage Docker build and deployment
- ✅ RBAC and K8s manifests
- ✅ Node details panel
- ✅ Connection status indicator

**Date Completed**: September 4, 2026

---

## Milestone 2: Interactive Features & Pluggable Layouts 🔄 IN PROGRESS

**Goal**: Add educational overlays, scenario-driven configuration, and flexible visualization system

### Completed in Milestone 2 ✅

**Logical Abstraction**:
- ✅ Scenario loader (ConfigMap + filesystem fallback)
- ✅ Logical graph mapper (K8s → architecture components)
- ✅ Scenario metadata models (patterns, technologies, layers, components, flows)
- ✅ API endpoints (`/api/scenarios`, `/api/scenario/{id}`, `/api/graph?scenario={id}`)

**Pluggable Layout System**:
- ✅ Layout configuration in `architecture.yaml`
- ✅ LayoutRouter component for layout selection
- ✅ VerticalLayout implementation (layered, top-to-bottom)
- ✅ HorizontalLayout placeholder
- ✅ Layout options (spacing, fitToViewport, compactCards)

**UI/UX Improvements**:
- ✅ LayerSidebar component (280px fixed width)
- ✅ Hover highlighting (layer → components)
- ✅ PatternOverlay for educational content
- ✅ ScenarioSwitcher dropdown in header
- ✅ Auto-scaling viewport (no scrolling)
- ✅ Health indicators (colored dots)
- ✅ Technology badges (e.g., "💡 Redpanda")
- ✅ Replica counts on components

**Scenario Configuration**:
- ✅ `architecture.yaml` schema definition
- ✅ Scenario-01 (Natural Disaster Supply Matching) with EDA + Saga patterns
- ✅ Pattern descriptions, benefits, and technologies
- ✅ Component definitions with K8s implementation mapping
- ✅ Flow definitions with edge styling

### Remaining in Milestone 2 📋

**Additional Layouts**:
- [ ] Implement HorizontalLayout (left-to-right pipeline flow)
  - For Pipes & Filters, Map-Reduce patterns
  - Horizontal Bezier curves for edges
  - Layer columns instead of rows

- [ ] Implement RadialLayout (hub-and-spoke)
  - For API Gateway, Service Mesh patterns
  - Center component positioning
  - Radial edge routing

- [ ] Implement ForceGraphLayout (physics-based network)
  - For complex microservices, Event Sourcing
  - Force simulation (attraction/repulsion)
  - Iterative layout convergence

- [ ] Implement GridLayout (matrix)
  - For Load Balancing, Scatter-Gather patterns
  - Configurable columns
  - Row/column headers

**Enhanced Interactions**:
- [ ] Click component to expand full details in NodeDetails panel
- [ ] Filter components by type (API, broker, worker, database)
- [ ] Search components by name
- [ ] Zoom/pan controls for large graphs
- [ ] Mini-map navigator
- [ ] Export graph as PNG/SVG

**Multi-Scenario Support**:
- [ ] Deploy scenario-02, scenario-03, etc.
- [ ] Side-by-side scenario comparison
- [ ] Scenario bookmarking
- [ ] Quick scenario switching without page reload

**Target Date**: End of September 2026

---

## Milestone 3: Live Data Flow Visualization 📋 PLANNED

**Goal**: Visualize real-time message flow through architecture with telemetry-driven animations

### Telemetry Collection

**Redpanda Metrics Collector**:
- [ ] Redpanda Admin API client in backend
- [ ] Poll topic metrics (message rates, consumer lag, producer count)
- [ ] Detect active producers/consumers
- [ ] Stream activity events via WebSocket

**PostgreSQL Metrics Collector**:
- [ ] Connection pool monitoring
- [ ] Query rate tracking
- [ ] Active connections count
- [ ] Transaction throughput

**HTTP Metrics Collector**:
- [ ] Request rate per endpoint
- [ ] Response time percentiles (p50, p95, p99)
- [ ] Error rate tracking
- [ ] Active connections

**WebSocket Event Types** (new):
```go
type DataFlowActivity struct {
    EdgeID    string  // Which edge is active
    Rate      float64 // Messages per second
    Direction string  // "produce" or "consume"
    Timestamp int64   // Unix timestamp
}
```

### Animated Visualizations

**Particle System**:
- [ ] SVG particle rendering in layout components
- [ ] Particles spawn at edge source, travel to target
- [ ] Speed/frequency based on message rate
- [ ] Color-coded by flow type (publish/subscribe/query)
- [ ] Opacity fade-in/fade-out

**Node Animations**:
- [ ] Pulse effect when actively processing
- [ ] Glow intensity based on throughput
- [ ] Health status transitions (smooth color changes)
- [ ] Replica count animations (scale up/down)

**Performance Overlays**:
- [ ] Real-time metrics display on hover
- [ ] Throughput graphs (sparklines)
- [ ] Latency heatmaps
- [ ] Alert indicators (warnings, errors)

### Data Flow Scenarios

**Example: Supply Offer Flow** (Scenario-01):
1. User submits supply offer → Supply API
2. API publishes event → Event Bus (particle animation)
3. Allocation Worker consumes event (node pulse)
4. Worker queries database (edge glow)
5. Worker publishes allocation → Event Bus (particle animation)

**Interactive Testing**:
- [ ] Test data generator UI
- [ ] "Send test message" buttons per scenario
- [ ] Playback controls (pause, speed up, slow down)
- [ ] Record and replay data flows

**Target Date**: Mid-October 2026

---

## Milestone 4: Advanced Scenarios & Patterns 📋 PLANNED

**Goal**: Deploy multiple scenarios covering different architecture patterns

### Additional Scenarios

**Scenario-02: Acoustic Wildlife Detection**:
- Patterns: Pipes & Filters → Transactional Outbox
- Layout: HorizontalLayout (data processing pipeline)
- Components: Audio Ingestion, FFT Processor, Classifier, Event Store

**Scenario-03: Surplus Food Redistribution**:
- Patterns: Pub-Sub → CQRS
- Layout: VerticalLayout with read/write separation
- Components: Command API, Query API, Event Bus, Write DB, Read DB

**Scenario-04: Environmental Sensor Network**:
- Patterns: Hexagonal Architecture → Event Sourcing
- Layout: RadialLayout (ports and adapters)
- Components: Core Domain, Adapters (HTTP, MQTT, gRPC), Event Store

**Scenario-05: Decentralized Community Health Network**:
- Patterns: Load Balancing + Circuit Breaker → Scatter-Gather
- Layout: GridLayout (replicated services)
- Components: Load Balancer, Service Instances, Circuit Breakers, Aggregator

### Pattern Library

**Pattern Documentation**:
- [ ] Detailed pattern descriptions with diagrams
- [ ] Benefits and trade-offs
- [ ] When to use / when not to use
- [ ] Implementation examples
- [ ] Common pitfalls

**Interactive Learning**:
- [ ] Step-by-step pattern walkthroughs
- [ ] Animated explanations (data flows, failure modes)
- [ ] Quizzes and knowledge checks
- [ ] Pattern comparison tool

### Developer Experience

**Scenario Authoring**:
- [ ] Scenario validation CLI tool (`cluster-viz validate scenario-01`)
- [ ] Hot-reload for `architecture.yaml` changes (watch mode)
- [ ] Visual scenario editor (drag-drop components)
- [ ] Component library browser
- [ ] Layout preview tool (render without deployment)

**Testing & Debugging**:
- [ ] Unit tests for mapper logic
- [ ] Integration tests for API endpoints
- [ ] E2E tests with Playwright
- [ ] Visual regression tests
- [ ] Performance benchmarks

**Target Date**: End of October 2026

---

## Milestone 5: Production Features 📋 FUTURE

**Goal**: Make cluster-viz production-ready for teaching and demonstrations

### Multi-Cluster Support

- [ ] Watch multiple clusters simultaneously
- [ ] Cluster switching UI
- [ ] Cross-cluster scenarios (multi-region, federation)
- [ ] Cluster comparison view

### Historical Data

- [ ] Time-travel through graph changes
- [ ] Playback historical data flows
- [ ] Incident timeline visualization
- [ ] Diff view (before/after comparisons)

### Collaboration

- [ ] Share scenarios via URL
- [ ] Embed graphs in teaching materials
- [ ] Export presentations (slide deck format)
- [ ] Annotation and commenting

### Observability Integration

- [ ] Prometheus metrics integration
- [ ] Jaeger trace visualization
- [ ] Log aggregation (Loki, Elasticsearch)
- [ ] Alert rule display

### Performance & Scalability

- [ ] Virtual scrolling for large graphs (1000+ nodes)
- [ ] Edge occlusion culling
- [ ] WebWorker for layout calculations
- [ ] Incremental graph updates (only changed nodes)
- [ ] Compressed WebSocket messages

**Target Date**: November 2026

---

## Success Metrics

### Milestone 2 (Current)
- ✅ Scenario-driven configuration working
- ✅ Pluggable layout system (1 implemented)
- ✅ Layer sidebar with hover highlighting
- ✅ Pattern overlay for education
- ⏳ Multiple layout types (target: 3+)
- ⏳ Scenario switching without page reload

### Milestone 3 (Upcoming)
- ⏳ Real-time telemetry collection (3 sources)
- ⏳ Particle animations on edges
- ⏳ Performance metrics overlay
- ⏳ Test data generator UI

### Milestone 4 (Future)
- ⏳ 5+ scenarios deployed
- ⏳ Pattern library with documentation
- ⏳ Scenario validation CLI
- ⏳ Visual scenario editor

### Milestone 5 (Future)
- ⏳ Multi-cluster support
- ⏳ Historical data time-travel
- ⏳ Collaboration features
- ⏳ Production-ready performance

---

## Contributing

When working on cluster-viz:

1. **Pick a milestone task** from the "Remaining" or "Planned" sections
2. **Create a feature branch**: `git checkout -b feature/horizontal-layout`
3. **Implement the feature** with tests
4. **Update this roadmap** with progress
5. **Commit with context**: Include milestone and feature in commit message
6. **Create a pull request** with description and screenshots

## Questions or Suggestions?

Open an issue or discussion in the repository to:
- Propose new milestones
- Suggest additional scenarios
- Request new layout types
- Report bugs or issues
