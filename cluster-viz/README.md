# Cluster-Viz: Architecture Pattern Visualizer

An interactive real-time dashboard that visualizes Kubernetes-based architecture patterns with educational overlays, pluggable layouts, and scenario-driven configurations.

## Features

### Core Visualization
- **Pluggable Layout System**: Multiple visualization types (vertical, horizontal, radial, grid, force-graph)
- **Scenario-Driven Configuration**: Each scenario defines its optimal layout via YAML
- **Layer Sidebar**: Discord-like sidebar showing architectural layers with hover highlighting
- **Logical Abstraction**: Maps physical K8s resources to logical architecture components
- **Real-time Updates**: WebSocket-based live updates from cluster state
- **Educational Overlays**: Pattern descriptions, benefits, and technology explanations

### Supported Layout Types

| Layout | Status | Best For |
|--------|--------|----------|
| **Vertical** | ✅ Implemented | Layered architectures, request/response flows |
| **Horizontal** | 📋 Planned | Linear pipelines, data processing |
| **Radial** | 📋 Planned | Hub-and-spoke, centralized components |
| **Force-Graph** | 📋 Planned | Complex interconnections, microservices |
| **Grid** | 📋 Planned | Matrix layouts, replicated services |

### Visual Features
- **Health Indicators**: Color-coded dots (green/yellow/red) for component health
- **Technology Badges**: Shows actual implementation (e.g., "💡 Redpanda", "💡 PostgreSQL")
- **Replica Counts**: Displays pod replica counts for scalable components
- **Animated Edges**: Pulsing animations for event flows
- **SVG Glow Effects**: Glowing edges with filters for visual depth
- **Dark Theme**: Navy/black gradient background with neon accents

## Tech Stack

**Backend**:
- Go + `client-go` (Kubernetes client)
- Scenario loader (ConfigMap + filesystem fallback)
- Logical graph mapper (K8s → logical components)
- WebSocket server (`nhooyr.io/websocket`)
- Embedded frontend (Go embed.FS)

**Frontend**:
- Svelte 4.2.20 (stable)
- Vite 5.4.21
- TypeScript
- Custom SVG-based rendering (no external graph libraries)
- Reactive stores for state management

**Deployment**:
- Multi-stage Docker build (Node → Go → Distroless)
- Single binary with embedded frontend
- Runs in k3d/Kind cluster

## Quick Start

### Prerequisites

- k3d or Kind cluster running
- Docker installed
- kubectl configured
- Helm (for infrastructure)

### Build and Deploy

```bash
# Build Docker image
cd cluster-viz
docker build -t cluster-viz:latest .

# Import to k3d cluster
k3d image import cluster-viz:latest -c cloud-architecture-sandbox

# Deploy to cluster
kubectl apply -f manifests/

# Wait for ready
kubectl wait --for=condition=Available deployment/cluster-viz -n cluster-viz --timeout=60s

# Port-forward to access UI
kubectl port-forward -n cluster-viz svc/cluster-viz 9090:80 &

# Open browser
open http://localhost:9090
```

### Deploy a Scenario

Each scenario provides an architecture.yaml file that defines:
- Metadata (title, description)
- Patterns (primary and secondary)
- Technologies used
- Layout configuration
- Logical layers
- Components (logical architecture)
- Flows (interactions)
- Telemetry collectors (for live data)

**Example: Deploy Scenario 01**

```bash
# Apply scenario ConfigMap
kubectl create configmap scenario-01-architecture \
  -n cluster-viz \
  --from-file=architecture.yaml=scenarios/scenario-01-eda/architecture.yaml \
  --dry-run=client -o yaml | kubectl apply -f -

# Label the ConfigMap
kubectl label configmap scenario-01-architecture -n cluster-viz \
  app=cluster-viz scenario=scenario-01

# Deploy scenario workloads (stubs or real implementations)
kubectl apply -f scenarios/scenario-01-eda/stubs.yaml
```

## Architecture

### Backend Structure

```
backend/
├── main.go                      # Entrypoint: clients, watcher, loader, handler, server
└── pkg/
    ├── models/
    │   ├── graph.go             # Physical graph (K8s resources)
    │   ├── scenario.go          # Scenario metadata + layout config
    │   └── events.go            # WebSocket event types
    ├── k8s/
    │   ├── client.go            # Typed + dynamic client setup
    │   ├── watcher.go           # Informers for K8s resources
    │   └── graph.go             # Physical graph builder
    ├── scenario/
    │   ├── loader.go            # Load from ConfigMap/filesystem
    │   └── mapper.go            # Map K8s → logical components
    └── api/
        └── handler.go           # REST + WebSocket + static files
```

### Frontend Structure

```
frontend/
└── src/
    ├── main.ts
    ├── App.svelte                      # Main app with sidebar + visualization layout
    ├── app.css
    └── lib/
        ├── types/
        │   ├── scenario.ts             # Scenario, LayoutConfig, LogicalGraph
        │   └── graph.ts                # Physical graph types (legacy)
        ├── stores/
        │   ├── scenario.ts             # Scenario state + fetching
        │   ├── graph.ts                # Physical graph state (legacy)
        │   └── websocket.ts            # WS connection with reconnect
        └── components/
            ├── LayerSidebar.svelte     # Left sidebar with layers
            ├── LayoutRouter.svelte     # Routes to layout engines
            ├── layouts/
            │   ├── VerticalLayout.svelte    # Top-to-bottom flow
            │   └── HorizontalLayout.svelte  # Pipeline (placeholder)
            ├── PatternOverlay.svelte   # Educational overlay
            ├── ScenarioSwitcher.svelte # Scenario dropdown
            ├── NodeDetails.svelte      # Sidebar panel for component details
            └── StatusBar.svelte        # Connection status (legacy)
```

### Data Flow

1. **Scenario Loading**:
   - Backend checks for ConfigMap `{scenario-id}-architecture` in cluster-viz namespace
   - Falls back to filesystem at `/scenarios/{scenario-id}-eda/architecture.yaml`
   - Parses YAML into `models.Scenario` including `Layout` field

2. **Logical Mapping**:
   - Backend watches physical K8s resources (Deployments, StatefulSets, Services, etc.)
   - Mapper reads `component.implementation.kubernetes` from scenario
   - Matches K8s resources to logical components
   - Builds logical graph with layers, components, edges, health

3. **Frontend Rendering**:
   - App fetches scenario metadata via `/api/scenario/{id}`
   - Fetches logical graph via `/api/graph?scenario={id}`
   - LayoutRouter reads `scenario.layout.type`
   - Routes to appropriate layout component (VerticalLayout, etc.)
   - Layout renders nodes, edges, hover highlighting

4. **Real-time Updates**:
   - WebSocket connection at `/api/ws`
   - Backend pushes ADDED/MODIFIED/DELETED events
   - Frontend updates graph reactively

## API Endpoints

### REST API

```bash
# List all scenarios
GET /api/scenarios
# Returns: [{ id, title, description, patterns[] }]

# Get full scenario metadata
GET /api/scenario/{scenario-id}
# Returns: { metadata, patterns, technologies, layout, layers, components, flows, telemetry }

# Get logical graph for scenario
GET /api/graph?scenario={scenario-id}
# Returns: { layers[], components{}, edges{}, health{} }
```

### WebSocket API

```bash
# Connect to live updates
WS /api/ws
# Receives: { type: "SYNC|ADDED|MODIFIED|DELETED", graph: {...} }
```

## Configuration

### Scenario Configuration (architecture.yaml)

```yaml
# Metadata
metadata:
  id: "scenario-01"
  title: "Natural Disaster Supply Matching"
  description: "..."

# Architecture patterns
patterns:
  primary:
    - name: "Event-Driven Architecture"
      description: "..."
      benefits: [...]
      color: "#f97316"
      icon: "⚡"
  secondary:
    - name: "Saga Pattern"
      ...

# Technology stack
technologies:
  - name: "Redpanda"
    role: "Event streaming platform"
    description: "..."

# Layout configuration
layout:
  type: "vertical"  # vertical | horizontal | radial | force-graph | grid
  options:
    spacing: "relaxed"
    fitToViewport: true
    compactCards: true

# Logical layers (top to bottom or left to right)
layers:
  - id: "presentation"
    name: "API Layer"
    description: "..."
    color: "#3b82f6"
    position: 0

# Logical components
components:
  - id: "supply-api"
    name: "Supply API"
    type: "api-gateway"
    layer: "presentation"
    role: "REST API Gateway"
    description: "..."
    icon: "🌐"
    implementation:
      kubernetes:
        deployment: "supply-api"
        service: "supply-api"
        namespace: "scenario-01"

# Component interactions
flows:
  - from: "supply-api"
    to: "event-bus"
    type: "publish"
    label: "Publishes SupplyOffered events"
    pattern: "Event-Driven Architecture"
    edgeStyle:
      color: "#f97316"
      style: "solid"
      animated: true
```

### Environment Variables

- `SCENARIOS_PATH`: Path to scenarios directory (default: `/scenarios`)
- `EXCLUDE_NAMESPACES`: Comma-separated list of namespaces to exclude from physical graph

### RBAC Permissions

The service account has read-only access to:
- Core: namespaces, pods, services, endpoints, configmaps
- Apps: deployments, statefulsets, daemonsets, replicasets
- Gateway API: gateways, httproutes

## Development

### Local Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Start dev server (proxies to backend)
npm run dev

# Build for production
npm run build
```

### Local Backend Development

```bash
cd backend

# Run locally (requires kubeconfig)
go run .

# Build binary
go build -o cluster-viz .

# Run tests
go test ./...
```

### Rebuild and Redeploy

```bash
# Rebuild Docker image
docker build -t cluster-viz:latest .

# Import to k3d
k3d image import cluster-viz:latest -c cloud-architecture-sandbox

# Restart deployment
kubectl rollout restart deployment/cluster-viz -n cluster-viz

# Watch rollout
kubectl rollout status deployment/cluster-viz -n cluster-viz
```

## Troubleshooting

### Pod not starting

```bash
kubectl describe pod -n cluster-viz -l app=cluster-viz
kubectl logs -n cluster-viz -l app=cluster-viz
```

### Scenario not loading

Check if ConfigMap exists:
```bash
kubectl get configmap -n cluster-viz -l app=cluster-viz
kubectl describe configmap scenario-01-architecture -n cluster-viz
```

Check backend logs:
```bash
kubectl logs -n cluster-viz -l app=cluster-viz | grep -i scenario
```

### Layout not rendering

1. Check scenario has layout field:
   ```bash
   curl http://localhost:9090/api/scenario/scenario-01 | jq '.layout'
   ```

2. Check browser console for errors:
   - Open DevTools → Console
   - Look for LayoutRouter or component errors

3. Verify logical graph has components:
   ```bash
   curl http://localhost:9090/api/graph?scenario=scenario-01 | jq '.components | length'
   ```

### WebSocket not connecting

Check that the service is running and port-forward is active:
```bash
kubectl get svc -n cluster-viz
kubectl port-forward -n cluster-viz svc/cluster-viz 9090:80
```

Check browser console for WebSocket errors.

## Cleanup

```bash
# Delete cluster-viz resources
kubectl delete -f manifests/

# Delete scenario ConfigMaps
kubectl delete configmap -n cluster-viz -l app=cluster-viz

# Delete namespace (removes everything)
kubectl delete namespace cluster-viz
```

## Future Enhancements

### Layouts
- [ ] Implement horizontal (pipeline) layout
- [ ] Implement radial (hub-and-spoke) layout
- [ ] Implement force-graph layout
- [ ] Implement grid (matrix) layout
- [ ] Custom positioning via layout config

### Live Data Flow
- [ ] Redpanda metrics collector (message rates, consumer lag)
- [ ] Particle animations along edges
- [ ] Real-time activity indicators
- [ ] Performance metrics visualization

### Interactive Features
- [ ] Click component to show full details
- [ ] Expand/collapse layers
- [ ] Filter components by type
- [ ] Search components by name
- [ ] Export graph as PNG/SVG

### Multi-Scenario
- [ ] Compare multiple scenarios side-by-side
- [ ] Switch scenarios without page reload
- [ ] Scenario history/timeline
- [ ] Bookmark favorite scenarios

### Advanced
- [ ] Historical view (time-travel through changes)
- [ ] Pod logs viewer in details panel
- [ ] Metrics integration (CPU/memory on nodes)
- [ ] Alert/notification system
- [ ] Multi-cluster support
