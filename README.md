# Cloud Architecture Sandbox

A hands-on learning sandbox for mastering software architecture patterns through real-world scenarios deployed in a local Kubernetes environment.

## Overview

This repository provides a structured approach to learning cloud-native architecture patterns by building, deploying, and validating incremental scenarios inside a local **Kind** (Kubernetes in Docker) cluster. Each scenario represents a real-world humanitarian or environmental use case, making the learning practical and meaningful.

## Learning Philosophy

Rather than abstract examples, this sandbox uses **real world project scenarios** to teach architecture patterns:

- 🏥 Healthcare networks for underserved populations
- 🌍 Environmental monitoring and climate research
- 🆘 Emergency response and disaster relief systems
- 🍲 Food redistribution and community aid

Each scenario follows a **two-phase learning approach**:

1. **Conceptual Teaching**: Understand the pattern through analogies, trade-offs, and minimal examples
2. **Practical Implementation**: Deploy a complete working system in your local Kind cluster

## Architecture Patterns Covered

### Scalability & Performance
- Load Balancing
- Scatter-Gather
- Map-Reduce
- Pipes and Filters / Stream Processing

### Resilience & Reliability
- Circuit Breaker
- Retry with Exponential Backoff
- Bulkhead

### Transaction & Data Patterns
- Event-Driven Architecture (EDA)
- Saga Pattern
- Transactional Outbox
- CQRS (Command Query Responsibility Segregation)
- Event Sourcing

### Integration & Extensibility
- Sidecar & Ambassador
- Anti-Corruption Layer
- Backends for Frontends (BFF)
- Hexagonal Architecture (Ports & Adapters)

### Communication
- Publisher-Subscriber with WebSockets

### Deployment
- Blue-Green Deployment
- Canary Deployment

## Learning Scenarios

1. **Natural Disaster Supply Matching System** - EDA → Saga Pattern
2. **Acoustic Wildlife & Threat Detection Engine** - Pipes & Filters → Transactional Outbox
3. **Surplus Food Redistribution Network** - Pub-Sub → CQRS
4. **Crowdsourced Environmental Sensor Network** - Hexagonal Architecture → Event Sourcing
5. **Decentralized Community Health Network** - Load Balancing + Circuit Breaker → Scatter-Gather
6. **Open Scientific Data Processing Platform** - Map-Reduce → Sidecar & Ambassador
7. **Multi-Language Humanitarian Aid Portal** - Backends for Frontends → Anti-Corruption Layer
8. **Progressive Emergency Alert System** - Blue-Green → Canary + Bulkhead

See [`learning-scenarios.md`](learning-scenarios.md) for detailed scenario descriptions.

## Tech Stack

### Core Infrastructure
- **Cluster**: Kind (Kubernetes in Docker)
- **Package Manager**: Helm (for all infrastructure deployments)
- **Ingress & Traffic Management**: Envoy Gateway using Kubernetes Gateway API
- **Messaging**: NATS, RabbitMQ, Redpanda, or Redis Pub/Sub (selected based on scenario requirements)
- **Data Stores**: PostgreSQL, Redis, InfluxDB, TimescaleDB (deployed per scenario)

### Design Principles
- All infrastructure components deployed via **Helm charts**
- All ingress/routing uses **Envoy Gateway** with Gateway API resources (`Gateway`, `HTTPRoute`, `ClientTrafficPolicy`, `BackendTrafficPolicy`)
- Messaging technology is selected based on pattern requirements:
  - **NATS**: IoT/edge computing, request/reply, scatter-gather, lightweight pub/sub
  - **RabbitMQ**: Work queues, enterprise integration, priority queues, reliable delivery
  - **Redpanda**: Event sourcing, saga patterns, high-throughput event logs (Kafka-compatible)
  - **Redis Pub/Sub**: Simple real-time notifications, ephemeral messaging

## Getting Started

### Prerequisites
- Docker Desktop or Docker Engine
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
- kubectl
- [Helm](https://helm.sh/docs/intro/install/) (for infrastructure deployments)

### Quick Start

```bash
# Clone the repository
git clone https://github.com/tymondragon/cloud-architecture-sandbox.git
cd cloud-architecture-sandbox

# Create the Kind cluster
kind create cluster --config cluster/kind-config.yaml

# Add Helm repositories
helm repo add envoy-gateway https://gateway.envoyproxy.io
helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm repo add rabbitmq https://charts.rabbitmq.com/
helm repo add redpanda https://charts.redpanda.com/
helm repo update

# Install base infrastructure (example - specific scenarios will detail requirements)
# Envoy Gateway
helm install eg envoy-gateway/gateway-helm --namespace envoy-gateway-system --create-namespace

# Choose a scenario and follow its implementation guide
# See learning-scenarios.md for details
```

## Repository Structure

```
.
├── README.md                   # This file
├── CLAUDE.md                   # AI assistant guidance for working in this repo
├── PROJECT_CONTEXT.md          # Repository overview and learning rules
├── architecture-context.md     # Teaching framework template
├── learning-scenarios.md       # 8 scenario domain catalog
├── cluster/                    # Kind cluster configurations
├── infra/                      # Shared infrastructure (Helm values files)
│   ├── envoy-gateway/          # Envoy Gateway Helm values & Gateway API resources
│   ├── nats/                   # NATS Helm values
│   ├── rabbitmq/               # RabbitMQ Helm values
│   ├── redpanda/               # Redpanda Helm values
│   └── databases/              # Database Helm values (PostgreSQL, Redis, etc.)
└── scenarios/                  # Per-scenario workloads
    ├── scenario-01-eda/
    ├── scenario-02-saga/
    └── ...
```

## Development Workflow

1. **Start with a scenario** from `learning-scenarios.md`
2. **Learn the pattern** using the teaching framework in `architecture-context.md`
3. **Deploy incrementally** - each scenario builds on previous infrastructure
4. **Validate** using provided commands (curl, kubectl logs, rpk)
5. **Clean up scenario workloads** while preserving shared infrastructure

## Contributing

This is a personal learning repository, but suggestions and improvements are welcome! Feel free to:

- Open issues for scenario ideas or corrections
- Submit PRs for improved implementations
- Share your learning experiences

## License

MIT License - feel free to use this for your own learning journey.

## Acknowledgments

- Architecture patterns inspired by [Michael Pogrebinsky's Complete Cloud Computing Software Architecture Patterns](https://www.udemy.com/course/the-complete-cloud-computing-software-architecture-patterns/)
- Built with guidance from Claude Code (Anthropic)

---

**Note**: This repository is under active development. Infrastructure manifests and scenario implementations will be added incrementally as patterns are learned and validated.
