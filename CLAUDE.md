# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Purpose

This is a **hands-on learning sandbox** for mastering software architecture patterns by building, deploying, and validating incremental scenarios in a local Kind (Kubernetes in Docker) cluster.

## Core Infrastructure Stack

- **Cluster**: Kind (Kubernetes in Docker)
- **Ingress & Traffic**: Envoy Gateway using Kubernetes Gateway API (`Gateway`, `HTTPRoute`, `ClientTrafficPolicy`, `BackendTrafficPolicy`)
- **Messaging**: NATS, RabbitMQ, Redpanda, or Redis Pub/Sub (selected based on scenario requirements)
- **Data Stores**: PostgreSQL, Redis, InfluxDB/TimescaleDB (deployed per scenario)
- **Package Manager**: Helm for all infrastructure deployments

**Non-Negotiable**:
- Always use Envoy Gateway for ingress/routing
- Always use Helm charts for deploying infrastructure components

**Messaging Selection**: Choose the messaging solution that best fits the pattern requirements:
- **NATS**: IoT/edge computing, request/reply patterns, scatter-gather, lightweight pub/sub
- **RabbitMQ**: Work queues, enterprise integration, priority queues, reliable delivery
- **Redpanda**: Event sourcing, saga patterns, high-throughput event logs (Kafka-compatible)
- **Redis Pub/Sub**: Simple real-time notifications, ephemeral messaging

## Learning Workflow (Two-Phase Approach)

### Phase 1: Conceptual Teaching
When teaching a new pattern:
1. **Core Concept**: Explain in 3 sentences with a real-world analogy
2. **Problem vs. Solution**: What problem it solves and what trade-offs it introduces
3. **Minimal Example**: Small code snippet showing component interaction
4. **Socratic Check**: Ask clarifying questions before moving to implementation

### Phase 2: Blueprint & Implementation
After concept confirmation, provide:
1. **Business Scenario**: Concrete use case from `learning-scenarios.md` with inputs/outputs
2. **Infrastructure Layer**:
   - Envoy Gateway routes and policies
   - Redpanda topics and configurations
   - Required databases and state stores
3. **Incremental Layering**: Explain how this builds on existing cluster state without tearing down shared infrastructure
4. **Component Breakdown**: List services/apps with specific responsibilities and business logic requirements
5. **Observability**: Port-forwards or HTTPRoute hosts for local access, plus logging/metrics commands
6. **Validation Steps**: Exact commands (curl, rpk, kubectl logs) to trigger and verify the pattern
7. **Clean Teardown**: Precise `kubectl delete` commands for scenario-specific workloads while preserving shared `infra/` components

## Repository Structure

```
.
├── architecture-context.md      # Teaching framework template
├── learning-scenarios.md        # Business domain catalog (4 scenarios)
├── PROJECT_CONTEXT.md          # Repository overview
├── cluster/                    # Kind cluster configs
├── infra/                      # Shared infrastructure manifests
│   ├── envoy-gateway/          # Gateway API resources
│   ├── redpanda/               # Event broker
│   └── storage/                # Database PVCs & StatefulSets
└── scenarios/                  # Per-scenario workloads
    ├── scenario-01-eda/
    ├── scenario-02-saga/
    └── ...
```

## Available Learning Scenarios

1. **Natural Disaster Supply Matching** (EDA → Saga Pattern)
2. **Acoustic Wildlife Detection** (Pipes & Filters → Transactional Outbox)
3. **Surplus Food Redistribution** (Pub-Sub → CQRS)
4. **Environmental Sensor Network** (Hexagonal Architecture → Event Sourcing)
5. **Decentralized Community Health Network** (Load Balancing + Circuit Breaker → Scatter-Gather)
6. **Open Scientific Data Processing** (Map-Reduce → Sidecar & Ambassador)
7. **Multi-Language Humanitarian Aid Portal** (Backends for Frontends → Anti-Corruption Layer)
8. **Progressive Emergency Alert System** (Blue-Green Deployment → Canary + Bulkhead)

See `learning-scenarios.md` for full details.

## Architecture Patterns Coverage

**Scalability & Performance:**
- Load Balancing
- Scatter-Gather
- Map-Reduce
- Pipes and Filters / Stream Processing

**Resilience & Reliability:**
- Circuit Breaker
- Retry with Exponential Backoff
- Bulkhead

**Transaction & Data Patterns:**
- Event-Driven Architecture (EDA)
- Saga Pattern (distributed transactions)
- Transactional Outbox
- CQRS (Command Query Responsibility Segregation)
- Event Sourcing

**Integration & Extensibility:**
- Sidecar & Ambassador
- Anti-Corruption Layer
- Backends for Frontends (BFF)
- Hexagonal Architecture (Ports & Adapters)

**Communication:**
- Publisher-Subscriber with WebSockets

**Deployment:**
- Blue-Green Deployment
- Canary Deployment

## Development Principles

- **Incremental**: Each scenario builds on previous infrastructure without full teardowns
- **Isolated Cleanup**: Scenario workloads can be removed independently while preserving shared infrastructure
- **Observable**: Every implementation includes verification commands and logging
- **Practical**: Focus on runnable, deployable code over theoretical examples
- **Interactive**: Use checking questions and step-by-step progression, not massive upfront implementations
- **Helm-First**: Use Helm charts for all infrastructure deployments (Envoy Gateway, messaging brokers, databases)

## Interactive Learning Skill

Use the `/learn-pattern` skill to start guided learning sessions:

```bash
# List all available scenarios
/learn-pattern

# Start a specific scenario (1-8)
/learn-pattern 1

# Start by pattern name
/learn-pattern eda
/learn-pattern "Saga Pattern"
```

The skill follows the two-phase approach: conceptual teaching first, then practical implementation with validation commands.

## Common Commands

### Cluster Management

```bash
# Create Kind cluster (if not exists)
kind create cluster --config cluster/kind-config.yaml

# Verify cluster
kubectl cluster-info --context kind-kind
kubectl get nodes

# Delete cluster (removes everything)
kind delete cluster
```

### Helm Repository Setup

```bash
# Add infrastructure Helm repositories
helm repo add envoy-gateway https://gateway.envoyproxy.io
helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm repo add rabbitmq https://charts.rabbitmq.com/
helm repo add redpanda https://charts.redpanda.com/
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# List added repositories
helm repo list
```

### Helm Deployment Workflow

```bash
# Install infrastructure component
helm install <release-name> <repo>/<chart> \
  --namespace <namespace> \
  --create-namespace \
  -f infra/<component>/values.yaml

# Upgrade with new values
helm upgrade <release-name> <repo>/<chart> \
  -f infra/<component>/values.yaml

# List installed releases
helm list --all-namespaces

# Uninstall (remove component)
helm uninstall <release-name> -n <namespace>

# Dry-run to preview changes
helm install --dry-run --debug <release-name> <repo>/<chart> -f values.yaml
```

### kubectl Operations

```bash
# Check deployments and pods
kubectl get pods -n <namespace>
kubectl get deployments -n <namespace>
kubectl get svc -n <namespace>

# Watch resources in real-time
kubectl get pods -n <namespace> -w

# Describe resources for debugging
kubectl describe pod <pod-name> -n <namespace>
kubectl describe svc <service-name> -n <namespace>

# View logs
kubectl logs -n <namespace> <pod-name> -f
kubectl logs -n <namespace> <pod-name> --previous  # Previous container logs

# Port-forward for local access
kubectl port-forward -n <namespace> svc/<service> 8080:80

# Apply manifests
kubectl apply -f scenarios/scenario-01/
kubectl apply -f infra/envoy-gateway/gateway.yaml

# Delete scenario workloads
kubectl delete -f scenarios/scenario-01/
```

### Messaging CLI Tools

```bash
# NATS
nats pub <subject> <message>
nats sub <subject>
nats stream list
nats consumer list <stream>

# RabbitMQ (exec into pod)
kubectl exec -n rabbitmq rabbitmq-0 -- rabbitmqctl list_queues
kubectl exec -n rabbitmq rabbitmq-0 -- rabbitmqctl list_exchanges

# Redpanda (rpk)
kubectl exec -n redpanda redpanda-0 -- rpk topic list
kubectl exec -n redpanda redpanda-0 -- rpk topic create <topic-name>
kubectl exec -n redpanda redpanda-0 -- rpk topic consume <topic-name>
kubectl exec -n redpanda redpanda-0 -- rpk topic produce <topic-name>

# Redis (redis-cli)
kubectl exec -n redis redis-0 -- redis-cli PING
kubectl exec -n redis redis-0 -- redis-cli KEYS '*'
kubectl exec -n redis redis-0 -- redis-cli SUBSCRIBE <channel>
kubectl exec -n redis redis-0 -- redis-cli PUBLISH <channel> <message>
```

### Gateway API Resources

```bash
# View Gateway API resources
kubectl get gateway -n envoy-gateway-system
kubectl get httproute -n <namespace>
kubectl get clienttrafficpolicy -n <namespace>
kubectl get backendtrafficpolicy -n <namespace>

# Describe for debugging
kubectl describe gateway <gateway-name> -n envoy-gateway-system
kubectl describe httproute <route-name> -n <namespace>
```

## Scenario Workflow

When working through a scenario:

1. **Read scenario details**: Check `learning-scenarios.md` for business context
2. **Prepare cluster**: Ensure Kind cluster exists and Helm repos are added
3. **Deploy shared infrastructure**: Envoy Gateway, messaging broker, databases (Helm)
4. **Create scenario manifests**: Place in `scenarios/scenario-XX/`
5. **Save Helm values**: Place in `infra/<component>/values.yaml` or `infra/<component>/scenario-XX-values.yaml`
6. **Apply and validate**: Use kubectl apply, check logs, test with curl/messaging tools
7. **Document validation commands**: Keep them in scenario README or implementation notes
8. **Clean teardown**: Delete scenario workloads, preserve shared infrastructure

## File Organization

- **Helm values files**: `infra/<component>/values.yaml` or `infra/<component>/scenario-specific.yaml`
- **Gateway API resources**: `infra/envoy-gateway/*.yaml` (shared) or `scenarios/scenario-XX/*.yaml` (scenario-specific)
- **Application manifests**: `scenarios/scenario-XX/` (Deployments, Services, ConfigMaps)
- **Cluster config**: `cluster/kind-config.yaml`
- **Documentation**: Each scenario should include validation commands and cleanup steps

## Troubleshooting

```bash
# Pod not starting
kubectl describe pod <pod-name> -n <namespace>  # Check events
kubectl logs <pod-name> -n <namespace>          # Check container logs

# Service not accessible
kubectl get svc -n <namespace>                  # Verify service exists
kubectl get endpoints -n <namespace>            # Check if pods are backing service

# Envoy Gateway routes not working
kubectl get httproute -n <namespace> -o yaml    # Check route configuration
kubectl logs -n envoy-gateway-system -l app.kubernetes.io/name=envoy-gateway -f

# Helm release issues
helm list --all-namespaces                      # Check release status
helm status <release-name> -n <namespace>       # View release details
helm rollback <release-name> -n <namespace>     # Rollback to previous version
```
