# Repository Context: Software Architecture Sandbox

## Project Objective
A hands-on learning repository designed to master advanced software architecture patterns by building, deploying, and validating incremental scenarios inside a local Kubernetes in Docker (**Kind**) cluster.

---

## Infrastructure Baseline
- **Local Cluster**: Kind (`cluster/kind-config.yaml`)
- **Ingress & Traffic Control**: **Envoy Gateway** using Kubernetes Gateway API (`Gateway`, `HTTPRoute`, `ClientTrafficPolicy`, `BackendTrafficPolicy`)
- **Event Streaming / Broker**: **Redpanda** (Kafka-API compatible, zero-JVM)
- **Data Stores**: PostgreSQL, Redis, InfluxDB/TimescaleDB (deployed on-demand per scenario)

---

## Repository Structure
.
├── claude-architecture-context.md   # Detailed teaching prompt framework
├── learning-scenarios.md            # Scenario domain catalog
├── cluster/
│   └── kind-config.yaml             # Cluster configuration & node port mappings
├── infra/
│   ├── envoy-gateway/               # Gateway API & Envoy manifests
│   ├── redpanda/                    # Redpanda cluster manifests
│   └── storage/                     # Database PVCs & StatefulSets
└── scenarios/
    ├── scenario-01-eda/             # Workloads for Event-Driven Architecture
    └── scenario-02-saga/            # Workloads for Saga Pattern

---

## Interactive Learning Rules
1. **Two-Phase Workflow**:
   - **Phase 1 (Concept)**: Explain the pattern (3-sentence analogy, problem/trade-offs, minimal code snippet, Socratic check).
   - **Phase 2 (Blueprint)**: Provide Kind cluster specs, Envoy Gateway routes, Redpanda topics, and validation commands.
2. **Incremental Scenarios**: Every pattern builds on top of previous infrastructure without tearing down core databases or brokers.
3. **Clean Teardowns**: Always provide exact `kubectl delete` commands for scenario-specific workloads while keeping shared `infra/` intact.
4. **Mandatory Tech**: Always enforce **Envoy Gateway** for traffic routing and **Redpanda** for messaging/events.