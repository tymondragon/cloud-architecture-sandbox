# Scenario 01: Event-Driven Architecture

## Quick Reference

**Pattern:** Event-Driven Architecture (EDA)
**Business Domain:** Natural Disaster Supply Matching System
**Tech Stack:** Redpanda, PostgreSQL, Envoy Gateway

## Quick Start

### Prerequisites

```bash
# Ensure Kind cluster exists
kind create cluster --config cluster/kind-config.yaml

# Add Helm repos
helm repo add envoy-gateway https://gateway.envoyproxy.io
helm repo add redpanda https://charts.redpanda.com/
helm repo add cnpg https://cloudnative-pg.io/charts/
helm repo update
```

### Deploy Infrastructure

```bash
# Envoy Gateway
helm install eg envoy-gateway/gateway-helm \
  --version v1.2.4 \
  --namespace envoy-gateway-system \
  --create-namespace \
  --wait

# Redpanda
helm install redpanda redpanda/redpanda \
  --namespace redpanda \
  --create-namespace \
  -f ../../infra/redpanda/values.yaml \
  --wait

# CloudNativePG Operator
helm install cnpg cnpg/cloudnative-pg \
  --namespace cnpg-system \
  --create-namespace \
  --wait

# PostgreSQL Cluster
kubectl apply -f ../../infra/postgresql/cluster.yaml
kubectl wait --for=condition=Ready cluster/supply-db -n databases --timeout=5m
```

### Deploy Application

```bash
# Apply Gateway and HTTPRoute
kubectl apply -f ../../infra/envoy-gateway/gateway.yaml
kubectl apply -f manifests/

# Create Redpanda topic
kubectl exec -n redpanda redpanda-0 -- \
  rpk topic create supply.requests --partitions 3 --replicas 1
```

### Test

```bash
# Port-forward Envoy Gateway
kubectl port-forward -n envoy-gateway-system \
  service/envoy-main-gateway 8080:80 &

# Send test request
curl -X POST http://localhost:8080/api/v1/supplies \
  -H "Host: supply-api.local" \
  -H "Content-Type: application/json" \
  -d '{
    "request_id": "TEST-001",
    "location": "shelter-downtown",
    "items": [{"type": "water_bottles", "quantity": 500}],
    "priority": "high"
  }'

# Watch logs
kubectl logs -n scenario-01 -l app=allocation-worker -f
```

### Cleanup

```bash
# Remove application only
kubectl delete namespace scenario-01

# Full cleanup (includes infrastructure)
kubectl delete -f ../../infra/postgresql/cluster.yaml
helm uninstall redpanda -n redpanda
helm uninstall cnpg -n cnpg-system
helm uninstall eg -n envoy-gateway-system
```

## Key Commands

### Monitoring

```bash
# Check pods
kubectl get pods -n scenario-01

# View API logs
kubectl logs -n scenario-01 -l app=supply-api -f

# View worker logs
kubectl logs -n scenario-01 -l app=allocation-worker -f

# Check Redpanda topics
kubectl exec -n redpanda redpanda-0 -- rpk topic list

# Check consumer lag
kubectl exec -n redpanda redpanda-0 -- rpk group describe allocation-workers
```

### Scaling

```bash
# Scale API
kubectl scale deployment -n scenario-01 supply-api --replicas=5

# Scale workers
kubectl scale deployment -n scenario-01 allocation-worker --replicas=10
```

## Full Lesson

See [lesson.md](./lesson.md) for complete conceptual teaching and implementation details.

## Next Pattern

Continue with **Saga Pattern** for distributed transaction management.
