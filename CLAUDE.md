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
