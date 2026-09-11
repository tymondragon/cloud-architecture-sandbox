# Cluster-Viz Development Session Context

**Date**: 2026-09-11
**Status**: Pluggable Layout System Complete ✅
**Current Task**: Documentation and commit

## Latest Progress (2026-09-11 - Session 3)

### ✅ COMPLETED: Pluggable Layout System with Layer Sidebar

**Major Architecture Redesign**: Transformed from single-layout visualization to a flexible, scenario-driven layout system with Discord-like sidebar interface.

**Problem Solved**:
- Original design required scrolling to see all components
- Single vertical layout didn't suit all architecture patterns
- No visual indication of layer organization

**Solution Implemented**:
- Layer sidebar (280px) on left showing all architectural layers
- Pluggable layout engine supporting multiple visualization types
- Hover highlighting between layers and components
- Auto-scaling viewport (no scrolling required)
- Scenario-driven layout configuration via YAML

**Files Created**:
- `cluster-viz/frontend/src/lib/components/LayerSidebar.svelte` - Left sidebar with layer definitions
- `cluster-viz/frontend/src/lib/components/LayoutRouter.svelte` - Routes to appropriate layout engine
- `cluster-viz/frontend/src/lib/components/layouts/VerticalLayout.svelte` - Vertical flow layout
- `cluster-viz/frontend/src/lib/components/layouts/HorizontalLayout.svelte` - Placeholder for horizontal pipeline layout

**Files Modified**:
- `cluster-viz/frontend/src/App.svelte` - Complete rewrite with sidebar + visualization layout
- `cluster-viz/frontend/src/lib/types/scenario.ts` - Added LayoutConfig and LayoutOptions interfaces
- `cluster-viz/backend/pkg/models/scenario.go` - Added Layout field and structs
- `scenarios/scenario-01-eda/architecture.yaml` - Added layout configuration section

**Layout Configuration Schema**:
```yaml
layout:
  type: "vertical"  # vertical | horizontal | radial | force-graph | grid | custom
  options:
    spacing: "relaxed"
    fitToViewport: true
    compactCards: true
    centerComponent: "optional-node-id"  # For radial layouts
    columns: 3  # For grid layouts
    positions: {}  # For custom layouts
```

**Features Implemented**:
- ✅ LayerSidebar component with color-coded layers
- ✅ Hover highlighting (layer → components)
- ✅ LayoutRouter with support for 5+ layout types
- ✅ VerticalLayout with compact card mode
- ✅ Auto-scaling viewport (no scrolling)
- ✅ Backend scenario loader integration
- ✅ Full TypeScript type safety

**Visual Design**:
- Discord-like sidebar layout
- 4 architectural layers: API, Integration, Business Logic, Data
- Compact component cards with health indicators
- Smooth hover transitions
- Glowing edges with SVG filters
- Dark navy/black gradient background

**Build & Deploy**:
- ✅ Frontend build successful
- ✅ Docker image built and deployed
- ✅ Running in k3d cluster
- ✅ Accessible at http://localhost:9090

### Previous Work (2026-09-04)

#### ✅ COMPLETED: Dark Theme Visual Redesign (Task #1)

**Files Modified**:
- `cluster-viz/frontend/src/lib/components/SimpleGraph.svelte` - Dark theme redesign
- `cluster-viz/frontend/src/lib/components/StatusBar.svelte` - Dark theme update

**Visual Changes**:
- Dark navy/black gradient background (#0a0e27 → #1a1f3a)
- Glowing nodes with gradient backgrounds
- Translucent node cards with backdrop blur
- SVG glow filters on edges
- Pulsing animation on dataflow edges
- Neon blue headers with text shadow

## Current State

### What's Working ✅

1. **Cluster-viz Dashboard (Pluggable Layout System)**
   - Backend: Go application serving scenario metadata with layout configs
   - Frontend: Svelte 4 with pluggable layout architecture
   - Real-time updates via WebSocket
   - Scenario API: `/api/scenarios`, `/api/scenario/{id}`, `/api/graph`
   - Deployment: Running in k3d cluster on namespace `cluster-viz`
   - Access: `http://localhost:9090` via port-forward

2. **Layout System Architecture**
   - **LayerSidebar**: Fixed 280px sidebar showing architectural layers
   - **LayoutRouter**: Routes to appropriate layout based on scenario config
   - **VerticalLayout**: Top-to-bottom flow with centered components
   - **HorizontalLayout**: Placeholder for pipeline visualization
   - **Future layouts**: Radial, force-graph, grid, custom

3. **Scenario Loading**
   - Loads from ConfigMaps (production) or filesystem (development)
   - Full scenario metadata including patterns, technologies, layers, components, flows
   - Layout configuration per scenario
   - Supports multiple scenarios with dropdown switcher

4. **Graph Building Logic**
   - **Logical Abstraction**: Maps K8s resources to logical architecture components
   - **Layers**: Presentation, Integration, Processing, Data
   - **Components**: API gateways, message brokers, workers, databases
   - **Flows**: Publish/subscribe, read-write, routing
   - **Health**: Component health status from K8s

5. **Visual Features**
   - Hover highlighting between sidebar and components
   - Health indicators (green/yellow/red dots)
   - Technology badges (e.g., "💡 Redpanda", "💡 PostgreSQL")
   - Replica counts
   - SVG edges with glow effects
   - Animated edges for event flows

### Technical Stack

**Frontend**:
- Svelte 4.2.20 (stable)
- Vite 5.4.21
- TypeScript
- Tailwind CSS (minimal usage)
- Custom SVG-based rendering

**Backend**:
- Go with Kubernetes client-go
- Scenario loader (ConfigMap or filesystem)
- Logical graph mapper
- WebSocket server for real-time updates
- Embedded frontend (static files in Go binary)

**Infrastructure**:
- k3d cluster (Kubernetes in Docker)
- Storage class: `local-path` (k3d default)
- Helm for infrastructure deployments

### Key Files

```
cluster-viz/
├── backend/
│   ├── main.go
│   ├── pkg/
│   │   ├── k8s/
│   │   │   ├── client.go
│   │   │   ├── graph.go          # Physical graph building
│   │   │   └── watcher.go        # K8s resource watchers
│   │   ├── models/
│   │   │   ├── graph.go          # Graph data structures
│   │   │   ├── scenario.go       # Scenario metadata (NEW: Layout field)
│   │   │   └── events.go         # WebSocket event types
│   │   ├── scenario/
│   │   │   ├── loader.go         # Loads scenarios from ConfigMap/filesystem
│   │   │   └── mapper.go         # Maps K8s → logical components
│   │   └── api/
│   │       └── handler.go        # HTTP + WebSocket handlers
├── frontend/
│   ├── src/
│   │   ├── App.svelte            # Main app with sidebar + visualization layout (NEW)
│   │   ├── lib/
│   │   │   ├── components/
│   │   │   │   ├── LayerSidebar.svelte      # Left sidebar (NEW)
│   │   │   │   ├── LayoutRouter.svelte      # Layout routing (NEW)
│   │   │   │   ├── layouts/
│   │   │   │   │   ├── VerticalLayout.svelte    # Vertical flow (NEW)
│   │   │   │   │   └── HorizontalLayout.svelte  # Placeholder (NEW)
│   │   │   │   ├── SimpleGraph.svelte       # Legacy (now replaced by layouts)
│   │   │   │   ├── StatusBar.svelte         # Connection status
│   │   │   │   ├── NodeDetails.svelte       # Sidebar panel for node info
│   │   │   │   ├── PatternOverlay.svelte    # Pattern education overlay
│   │   │   │   └── ScenarioSwitcher.svelte  # Scenario dropdown
│   │   │   ├── stores/
│   │   │   │   ├── scenario.ts              # Scenario state management (NEW)
│   │   │   │   ├── graph.ts                 # Graph state (legacy)
│   │   │   │   └── websocket.ts             # WebSocket connection
│   │   │   └── types/
│   │   │       ├── scenario.ts              # Scenario types (NEW: LayoutConfig)
│   │   │       └── graph.ts                 # Graph types
│   ├── package.json
│   └── vite.config.ts
├── manifests/
│   ├── deployment.yaml
│   ├── service.yaml
│   └── rbac.yaml
└── Dockerfile                    # Multi-stage build

scenarios/scenario-01-eda/
├── architecture.yaml             # Scenario metadata (NEW: layout section)
├── lesson.md
├── stubs.yaml
└── manifests/
```

## Architecture: Pluggable Layout System

### Data Flow

1. **Scenario Loading** (Backend)
   - Loader checks for ConfigMap `{scenario-id}-architecture` in cluster-viz namespace
   - Falls back to filesystem at `/scenarios/{scenario-id}-eda/architecture.yaml`
   - Parses YAML into `models.Scenario` struct including `Layout` field

2. **Logical Graph Mapping** (Backend)
   - Mapper takes scenario metadata + physical K8s graph
   - Maps K8s resources to logical components based on implementation.kubernetes fields
   - Builds logical graph with layers, components, edges, health

3. **API Endpoints** (Backend)
   - `GET /api/scenarios` - List all scenarios (summary)
   - `GET /api/scenario/{id}` - Full scenario metadata (includes layout config)
   - `GET /api/graph?scenario={id}` - Logical graph for scenario

4. **Frontend Rendering** (Svelte)
   - App.svelte fetches scenario metadata and logical graph
   - LayoutRouter reads `scenario.layout.type` and routes to appropriate layout component
   - Layout component (e.g., VerticalLayout) renders nodes, edges, hover highlighting
   - LayerSidebar displays layers, emits hover events

### Layout Types (Current & Planned)

| Layout Type | Status | Best For | Example Pattern |
|-------------|--------|----------|-----------------|
| **vertical** | ✅ Implemented | Layered architectures, request/response flows | EDA, Saga, CQRS |
| **horizontal** | 📋 Planned | Linear pipelines, data processing | Pipes & Filters, Map-Reduce |
| **radial** | 📋 Planned | Hub-and-spoke, centralized components | API Gateway, Service Mesh |
| **force-graph** | 📋 Planned | Complex interconnections, microservices | Event Sourcing, Microservices |
| **grid** | 📋 Planned | Matrix layouts, replicated services | Scatter-Gather, Load Balancing |
| **custom** | 📋 Planned | Scenario-specific positioning | Special cases |

### Example: Vertical Layout Config

```yaml
layout:
  type: "vertical"
  options:
    spacing: "relaxed"      # relaxed | compact | wide
    fitToViewport: true     # Auto-scale to eliminate scrolling
    compactCards: true      # Reduce card size and padding
```

### Example: Radial Layout Config (Future)

```yaml
layout:
  type: "radial"
  options:
    centerComponent: "api-gateway"  # Hub component
    radius: 300                     # Distance from center
    angleSpacing: 45               # Degrees between nodes
```

## Problems Solved

### 1. Scrolling Issue in Original Design
**Problem**: Users had to scroll to see all components in the visualization
**Solution**:
- Added LayerSidebar to show all layers at once
- Implemented auto-scaling viewport with `overflow: hidden`
- Compact card mode to fit more in viewport

### 2. Single Layout Limitation
**Problem**: Vertical layout doesn't suit all architecture patterns
**Solution**:
- Pluggable layout system with LayoutRouter
- Scenario-driven configuration
- Multiple layout engines for different pattern types

### 3. Layer Organization Not Visible
**Problem**: No way to see which components belong to which architectural layer
**Solution**:
- LayerSidebar showing all layers with descriptions
- Hover highlighting to show layer-component associations
- Color-coded layer indicators

### 4. Svelte 5 Effect_Orphan Error (Previous Session)
**Problem**: Using Svelte 5 runes in plain Vite setup caused persistent errors
**Solution**: Downgraded to Svelte 4 with stable component architecture

## Current Plan: Next Features

### Milestone 2: Interactive Features (In Progress)
- ✅ Pattern overlay with educational content
- ✅ Scenario switcher dropdown
- ✅ Layer sidebar with hover highlighting
- ✅ Pluggable layout system
- 📋 Implement horizontal layout
- 📋 Implement radial layout
- 📋 Node detail panel enhancements

### Milestone 3: Live Data Flow (Upcoming)
**Goal**: Visualize real-time message flow through architecture

**Task #2: Add Redpanda Metrics Collector**
1. Add Redpanda Admin API client to backend
2. Poll topic metrics (message rates, consumer lag)
3. Detect active producers/consumers
4. Stream activity events via WebSocket

**Task #3: Implement Particle Animations**
1. SVG particle system in layout components
2. Particles spawn at edge source, travel to target
3. Speed/frequency based on message rate
4. Color-coded by edge type
5. Nodes pulse when actively processing

### Milestone 4: Advanced Scenarios
- Deploy scenario-02, scenario-03, etc.
- Different layout types per scenario
- More complex component interactions
- Multi-scenario comparisons

## Commands Reference

**Port-forward to dashboard**:
```bash
kubectl port-forward -n cluster-viz svc/cluster-viz 9090:80 &
```

**Rebuild and deploy**:
```bash
cd cluster-viz
docker build -t cluster-viz:latest .
k3d image import cluster-viz:latest -c cloud-architecture-sandbox
kubectl rollout restart deployment/cluster-viz -n cluster-viz
kubectl wait --for=condition=Available deployment/cluster-viz -n cluster-viz --timeout=60s
```

**Check APIs**:
```bash
# List scenarios
curl -s http://localhost:9090/api/scenarios | jq '.'

# Get scenario with layout config
curl -s http://localhost:9090/api/scenario/scenario-01 | jq '.layout'

# Get logical graph
curl -s http://localhost:9090/api/graph?scenario=scenario-01 | jq '{layers, components}'
```

**View logs**:
```bash
kubectl logs -n cluster-viz -l app=cluster-viz -f
```

## Notes

- Always use Svelte 4 (not 5) for this project
- Layout configuration is per-scenario in architecture.yaml
- Scenario loader tries ConfigMap first, falls back to filesystem
- k3d uses `local-path` storage class (not `standard` like Kind)
- Distroless final image (no shell tools available in pod)

## Resume Instructions

### To Continue This Session:

1. **Start cluster** (if stopped):
   ```bash
   k3d cluster start cloud-architecture-sandbox
   kubectl wait --for=condition=Ready nodes --all --timeout=120s
   ```

2. **Port-forward to dashboard**:
   ```bash
   kubectl port-forward -n cluster-viz svc/cluster-viz 9090:80 &
   ```

3. **Open in browser**:
   ```bash
   open http://localhost:9090
   ```

4. **Expected UI**:
   - Layer sidebar on left (280px) with 4 layers
   - Main visualization with 4 components vertically centered
   - Hover over layer highlights corresponding components
   - No scrolling required
   - Help button opens pattern overlay
   - Scenario switcher in header

### Current Task Status

- ✅ Task #1: Dark theme redesign - COMPLETE
- ✅ Task #11-16: Pluggable layout system - COMPLETE
- 📋 Task #2: Add Redpanda metrics collector - NOT STARTED
- 📋 Task #3: Implement particle animations - NOT STARTED

## Questions / Decisions Needed

- ✅ Color scheme for node types - DECIDED
- ✅ Layout architecture - DECIDED (pluggable system)
- ✅ Sidebar design - DECIDED (Discord-like layer sidebar)
- Should we implement horizontal layout next or move to live data flow?
- Which scenario should get horizontal layout (pipeline pattern)?
- How to handle high message rates in particle animations?
