# Cluster Visualization Dashboard

An interactive real-time dashboard that visualizes your Kubernetes cluster topology, showing resources, their relationships, and data flow.

## Features

- **Real-time updates** via WebSocket
- **Auto-discovery** of 9 resource types: Namespaces, Pods, Services, Deployments, StatefulSets, DaemonSets, Gateways, HTTPRoutes, and Endpoints
- **5 edge types** showing relationships:
  - **Ownership** (gray dashed): Controller → Pod relationships
  - **Routing** (blue solid): HTTPRoute → Service routing
  - **Selection** (gray dotted): Service → Pod label matching
  - **Data Flow** (orange animated): Pod → Service connections inferred from environment variables
  - **Parent Ref** (purple solid): HTTPRoute → Gateway references
- **Interactive diagram** with zoom, pan, and click-to-inspect
- **Status indicators** showing health of each resource
- **Namespace filtering** to exclude system namespaces

## Tech Stack

- **Backend**: Go + `client-go` (typed + dynamic) + `nhooyr.io/websocket`
- **Frontend**: Svelte 5 + Vite + Tailwind CSS + Svelte Flow (diagram) + Dagre (auto-layout)
- **Deployment**: Single Docker image with embedded SPA

## Quick Start

### Prerequisites

- Kind cluster running
- Docker installed
- kubectl configured
- Envoy Gateway deployed (for HTTPRoute exposure)

### Build and Deploy

```bash
# Build Docker image and load into Kind
cd scenarios/cluster-viz
bash scripts/build.sh

# Deploy to cluster
bash scripts/deploy.sh

# Port-forward to access UI
kubectl port-forward -n cluster-viz svc/cluster-viz 9090:80

# Open browser
open http://localhost:9090
```

## Verification Commands

### Check Deployment

```bash
# Check pod status
kubectl get pods -n cluster-viz

# Check logs
kubectl logs -n cluster-viz -l app=cluster-viz -f

# Check service
kubectl get svc -n cluster-viz
```

### Test REST API

```bash
# Get full graph snapshot
curl http://localhost:9090/api/graph | jq .

# Test WebSocket (requires websocat)
websocat ws://localhost:9090/api/ws
```

### Via Envoy Gateway

If you have Envoy Gateway configured with a main-gateway:

```bash
# Add to /etc/hosts
echo "127.0.0.1 cluster-viz.local" | sudo tee -a /etc/hosts

# Port-forward the gateway (if not already)
kubectl port-forward -n envoy-gateway-system svc/envoy-main-gateway 80:80

# Access via hostname
curl http://cluster-viz.local/api/graph | jq .
open http://cluster-viz.local
```

## Architecture

### Backend Structure

```
backend/
├── main.go                   # Entrypoint: wire clients, watcher, handler, serve
└── pkg/
    ├── models/
    │   ├── graph.go          # GraphNode, GraphEdge, Graph types
    │   └── events.go         # GraphEvent envelope (ADDED/MODIFIED/DELETED/SYNC)
    ├── k8s/
    │   ├── client.go         # Typed + dynamic client setup
    │   ├── watcher.go        # Informers for 9 resource types
    │   └── graph.go          # Graph builder with edge inference
    └── api/
        └── handler.go        # REST + WebSocket + static file serving
```

### Edge Inference Logic

1. **Ownership**: Walks `ownerReferences` chain, skips ReplicaSets, creates direct Deployment→Pod edges
2. **ParentRef**: HTTPRoute `spec.parentRefs` → Gateway
3. **Routing**: HTTPRoute `spec.rules[].backendRefs` → Service
4. **Selection**: Service `spec.selector` matched against Pod labels
5. **DataFlow**: Scans Pod container env vars for connection patterns:
   - `postgres://`, `nats://`, `amqp://`, `redis://`
   - `*.svc.cluster.local`
   - `REDPANDA_BROKERS`

### Frontend Structure

```
frontend/
└── src/
    ├── main.ts
    ├── App.svelte
    ├── app.css
    └── lib/
        ├── types/graph.ts          # TypeScript mirrors of Go models
        ├── stores/
        │   ├── graph.ts            # Reactive graph state + dagre layout
        │   └── websocket.ts        # WS connection with reconnect
        └── components/
            ├── FlowDiagram.svelte  # Main Svelte Flow diagram
            ├── KubeNode.svelte     # Custom node (status colors, kind icons)
            ├── StatusBar.svelte    # Connection status + legend
            └── NodeDetails.svelte  # Click-to-inspect panel
```

## Configuration

### Environment Variables

- `EXCLUDE_NAMESPACES`: Comma-separated list of namespaces to exclude (default: `kube-system,kube-public,kube-node-lease,local-path-storage`)
- `KUBECONFIG`: Path to kubeconfig file (for local dev, auto-detects `~/.kube/config`)

### RBAC Permissions

The service account has read-only access to:
- Core: namespaces, pods, services, endpoints
- Apps: deployments, statefulsets, daemonsets, replicasets
- Gateway API: gateways, httproutes

## Development

### Local Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Start dev server
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
```

### Rebuild and Redeploy

```bash
# Rebuild image
bash scripts/build.sh

# Restart deployment
kubectl rollout restart deployment/cluster-viz -n cluster-viz
```

## Troubleshooting

### Pod not starting

```bash
kubectl describe pod -n cluster-viz -l app=cluster-viz
kubectl logs -n cluster-viz -l app=cluster-viz
```

### WebSocket not connecting

Check that the service is running and port-forward is active:

```bash
kubectl get svc -n cluster-viz
kubectl port-forward -n cluster-viz svc/cluster-viz 9090:80
```

### No edges showing

Edges are computed based on metadata. Check that:
- Services have selectors
- Pods have ownerReferences
- HTTPRoutes have parentRefs and backendRefs
- Pods have environment variables with connection strings

### Gateway API resources not showing

Ensure Gateway API CRDs are installed:

```bash
kubectl get crds | grep gateway
```

If missing, install Envoy Gateway or Gateway API CRDs.

## Cleanup

```bash
# Delete cluster-viz resources
kubectl delete -f manifests/

# Delete namespace (removes everything)
kubectl delete namespace cluster-viz
```

## Future Enhancements

- [ ] Namespace grouping (parent nodes containing resources)
- [ ] Resource filtering by kind or namespace
- [ ] Search/highlight specific resources
- [ ] Historical view (time-travel through graph changes)
- [ ] Export graph as PNG/SVG
- [ ] Metrics integration (CPU/memory usage on nodes)
- [ ] Pod logs viewer in details panel
- [ ] ConfigMap/Secret visualization
- [ ] Ingress resource support
