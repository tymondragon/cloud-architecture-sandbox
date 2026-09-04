# Cluster-Viz Development Session Context

**Date**: 2026-09-04
**Status**: In Progress - Visual Redesign Phase

## Current State

### What's Working ✅

1. **Cluster-viz Dashboard (Basic Version)**
   - Backend: Go application that watches Kubernetes resources and builds a graph
   - Frontend: Svelte 4 with custom SimpleGraph component (no external graph libraries)
   - Real-time updates via WebSocket
   - API endpoint: `/api/graph` returns filtered graph data
   - Deployment: Running in k3d cluster on namespace `cluster-viz`
   - Access: `http://localhost:9090` via port-forward

2. **Graph Filtering**
   - Successfully filters out infrastructure namespaces (cert-manager, envoy-gateway-system, cnpg-system, kube-system)
   - Shows only lesson-relevant components (19 nodes, 18 edges currently)
   - Includes: scenario-* namespaces, redpanda, databases, Gateway API resources

3. **Graph Building Logic**
   - **Nodes**: Deployments, Pods, Services, StatefulSets, HTTPRoutes, Gateways
   - **Edges**:
     - `ownership`: Deployment → Pods (via ownerReferences)
     - `selection`: Service → Pods (via label selectors)
     - `routing`: HTTPRoute → Service (via backendRefs)
     - `parentRef`: Gateway → HTTPRoute (via parentRefs)
     - `dataflow`: Pod → Service (detected via env vars with connection strings)

4. **Test Infrastructure Deployed**
   - k3d cluster: `cloud-architecture-sandbox` (1 control-plane + 2 workers)
   - Envoy Gateway (for Gateway API)
   - Redpanda (event broker)
   - CloudNativePG operator + PostgreSQL cluster
   - Scenario-01 stubs: supply-api and allocation-worker (nginx placeholders)

### Technical Stack

**Frontend**:
- Svelte 4.2.20 (stable - Svelte 5 caused effect_orphan errors)
- Vite 5.4.21
- TypeScript
- Tailwind CSS
- Custom SVG-based graph rendering (no @xyflow/svelte)

**Backend**:
- Go with Kubernetes client-go
- Watches: Deployments, Pods, Services, StatefulSets, HTTPRoutes, Gateways
- WebSocket server for real-time graph updates
- Embedded frontend (static files served from Go binary)

**Infrastructure**:
- k3d cluster (Kubernetes in Docker)
- Storage class: `local-path` (k3d default)
- Helm used for all infrastructure deployments

### Key Files

```
cluster-viz/
├── backend/
│   ├── main.go
│   ├── pkg/
│   │   ├── k8s/
│   │   │   ├── client.go
│   │   │   ├── graph.go          # Graph building logic + filtering
│   │   │   └── watcher.go        # K8s resource watchers
│   │   ├── models/
│   │   │   ├── graph.go          # Graph data structures
│   │   │   └── events.go         # WebSocket event types
│   │   └── api/
│   │       └── handler.go        # HTTP + WebSocket handlers
├── frontend/
│   ├── src/
│   │   ├── App.svelte            # Main app (subscribes to stores, passes to SimpleGraph)
│   │   ├── lib/
│   │   │   ├── components/
│   │   │   │   ├── SimpleGraph.svelte   # Custom graph visualization
│   │   │   │   ├── StatusBar.svelte     # Connection status + stats
│   │   │   │   └── NodeDetails.svelte   # Sidebar panel for node info
│   │   │   ├── stores/
│   │   │   │   ├── graph.ts      # Svelte stores for nodes/edges
│   │   │   │   └── websocket.ts  # WebSocket connection management
│   │   │   └── types/
│   │   │       └── graph.ts      # TypeScript types
│   ├── package.json              # Svelte 4, Vite 5
│   └── vite.config.ts
├── manifests/
│   ├── deployment.yaml
│   ├── service.yaml
│   └── rbac.yaml
└── Dockerfile                    # Multi-stage build (frontend + backend)
```

## Problems We Solved

### 1. Svelte 5 Effect_Orphan Error
**Problem**: Using Svelte 5 runes in plain Vite setup caused persistent `effect_orphan` errors
**Root Cause**: IDP (reference project) uses SvelteKit, not plain Vite. SvelteKit handles component initialization differently.
**Solution**: Downgraded to Svelte 4 with the custom graph component

### 2. @xyflow/svelte Incompatibility
**Problem**: @xyflow/svelte library had Svelte 5 compatibility issues
**Solution**: Built custom SimpleGraph component based on IDP's InfraDiagram pattern (pure SVG + HTML)

### 3. Too Much Infrastructure Noise
**Problem**: Graph showed all cluster resources (50+ nodes) including cert-manager, operators, etc.
**Solution**: Added `isLessonRelevant()` filter in `backend/pkg/k8s/graph.go` to only show scenario and data infrastructure

## Current Plan: Visual Redesign + Live Data Flow

### Task #1: Redesign Graph with Dark Theme ⏳ IN PROGRESS
**Goal**: Transform boring white boxes into sleek, glowing visualization

**Visual Changes**:
- Dark navy/black background (#0a0e27 or similar)
- Glowing nodes with CSS gradients and box-shadow
- Neon accent colors:
  - Cyan/blue for standard nodes
  - Orange for dataflow edges (with glow)
  - Purple for parentRef edges
  - Blue for routing edges
- Smooth animations and transitions
- Better typography (light text on dark)
- Node icons/badges with glow effects

**Reference**: User provided image of dark-themed data pipeline diagram with glowing nodes and animated flow

**Files to Modify**:
- `frontend/src/lib/components/SimpleGraph.svelte` (complete redesign of styles)
- May need to update `App.svelte` for dark background

### Task #2: Add Redpanda Metrics Collector 📋 TODO
**Goal**: Monitor real-time message flow through Redpanda

**Implementation**:
1. Add Redpanda Admin API client to backend
2. Poll topic metrics (message rates, consumer lag)
3. Detect active producers/consumers
4. Stream activity events via WebSocket

**WebSocket Event Types** (new):
```go
type DataFlowActivity struct {
    EdgeID    string  // Which edge is active
    Rate      float64 // Messages per second
    Direction string  // "produce" or "consume"
}
```

**Files to Create/Modify**:
- `backend/pkg/telemetry/redpanda.go` (new - Redpanda client)
- `backend/pkg/models/events.go` (add DataFlowActivity type)
- `backend/pkg/api/handler.go` (broadcast activity events)

### Task #3: Implement Particle Animations 📋 TODO
**Goal**: Visualize live data flowing through edges

**Implementation**:
1. SVG particle system in SimpleGraph.svelte
2. Particles spawn at edge source, travel to target
3. Speed/frequency based on message rate from telemetry
4. Color-coded by edge type
5. Nodes pulse when actively processing

**CSS/Animation**:
- SVG `<circle>` elements animated along path
- `animateMotion` or JavaScript-based animation
- Glow filters for particles
- Opacity fade-in/fade-out

**Files to Modify**:
- `frontend/src/lib/components/SimpleGraph.svelte` (add particle rendering)
- `frontend/src/lib/stores/websocket.ts` (handle activity events)

## Next Steps

1. **Complete Visual Redesign** (Task #1)
   - Rewrite SimpleGraph.svelte with dark theme
   - Add glowing effects to nodes and edges
   - Test in browser

2. **Add Redpanda Telemetry** (Task #2)
   - Implement Redpanda Admin API client
   - Stream activity data via WebSocket
   - Test with scenario-01 stubs

3. **Implement Particle Animations** (Task #3)
   - Add SVG particle system
   - Animate based on live telemetry
   - Polish animations and timing

4. **Testing & Validation**
   - Generate real traffic (curl to supply-api, trigger Redpanda messages)
   - Verify particles animate along correct edges
   - Validate performance with high message rates

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

**Check graph API**:
```bash
curl -s http://localhost:9090/api/graph | python3 -m json.tool
```

**Test Redpanda (when implemented)**:
```bash
kubectl exec -n redpanda redpanda-0 -- rpk topic produce test-topic
kubectl exec -n redpanda redpanda-0 -- rpk topic consume test-topic
```

## Notes

- Always use Svelte 4 (not 5) for this project
- Custom graph component works better than @xyflow for our use case
- Filter infrastructure namespaces to keep graph focused on lessons
- Real-time updates via WebSocket are working well
- k3d uses `local-path` storage class (not `standard` like Kind)

## Questions / Decisions Needed

- What color scheme for different node types? (Deployment, Pod, Service, etc.)
- Should particles be persistent or fade out after traveling?
- How to handle high message rates (thousands per second)?
- Should we aggregate activity data or show individual messages?
