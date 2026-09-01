# Software Architecture Scenario Catalog

A collection of system scenarios designed for local execution inside a Kind (Kubernetes in Docker) cluster using Envoy Gateway for ingress/traffic control and Redpanda for high-throughput event streaming.

---

## Scenario 1: Natural Disaster Supply Matching System

### Business Domain

An emergency relief logistics system designed to match urgent supplies (water, medicine, shelters) with field requests in real time during power and network instability.

### Baseline Implementation

- **Architecture**: Monolithic REST API + Relational Database (PostgreSQL).
- **Behavior**: Direct synchronous write endpoints where field agents submit supply requests and donors post available resources.

### Architecture Evolution

- **Phase 1: Event-Driven Architecture (EDA)**
  - Decouple request ingestion from allocation algorithms using an asynchronous event broker (Redpanda/RabbitMQ).
  - API returns `202 Accepted` immediately; worker pods consume supply/demand events off the queue.
- **Phase 2: The Saga Pattern**
  - Implement distributed transaction management (Reserve Items -> Allocate Logistics -> Update Shelter Queue).
  - Provide compensating rollback logic if transit routes become blocked or unavailable.

### Kind & Envoy Infrastructure

- **Databases**: PostgreSQL (Relational state).
- **Messaging**: Redpanda or RabbitMQ.
- **Envoy Gateway Traffic Policies**: Envoy Gateway `ClientTrafficPolicy` for rate limiting on API ingestion routes (`HTTPRoute`) to prevent traffic spikes during emergency events from dropping requests.

---

## Scenario 2: Acoustic Wildlife & Threat Detection Engine

### Business Domain

An automated IoT acoustic monitoring network that ingests high-frequency audio telemetry from remote forest sensors to identify localized environmental threats (e.g., chainsaws, vehicle engines, gunshots).

### Baseline Implementation

- **Architecture**: Simple HTTP Data Collector.
- **Behavior**: Remote sensors post JSON metadata payloads containing frequency profiles to a single service endpoint that writes raw logs to a database.

### Architecture Evolution

- **Phase 1: Pipes and Filters / Stream Processing**
  - Stream raw sensor telemetry into Redpanda topics.
  - Pass streams through specialized processing pods:
    - *Filter 1*: Background noise reduction.
    - *Filter 2*: Audio fingerprinting & ML model evaluation.
    - *Filter 3*: Threat classifier & alerting event dispatch.
- **Phase 2: Transactional Outbox Pattern**
  - Ensure local nodes safely write alert states to a transactional outbox before forwarding over intermittent satellite links, preventing alert loss on network drops.

### Kind & Envoy Infrastructure

- **Databases**: TimescaleDB / InfluxDB (Time-series data).
- **Messaging**: Redpanda (Kafka API compatible, zero-JVM footprint).
- **Envoy Gateway Traffic Policies**: Envoy Gateway `BackendTrafficPolicy` for circuit breaking and timeout rules on ingress routes to gracefully handle high-latency or unstable satellite links.

---

## Scenario 3: Surplus Food Redistribution Network

### Business Domain

A real-time matching network connecting commercial food suppliers, distribution centers, and field delivery volunteers to redistribute perishable food before expiration.

### Baseline Implementation

- **Architecture**: Web application driven by standard HTTP polling.
- **Behavior**: Receivers query the primary database every few minutes to pull list updates for nearby food donations.

### Architecture Evolution

- **Phase 1: Publisher-Subscriber & Persistent WebSockets**
  - Replace HTTP polling with persistent WebSocket connections routed through Envoy Gateway.
  - Broadcast new donation notifications dynamically to driver instances based on location subscriptions using Redis Pub/Sub.
- **Phase 2: CQRS (Command Query Responsibility Segregation)**
  - Separate write-heavy donation intake workflows (PostgreSQL) from read-heavy geospatial and dietary search queries (Elasticsearch/Redis).

### Kind & Envoy Infrastructure

- **Databases**: PostgreSQL (Writes) + Redis (Read-Cache / Pub-Sub).
- **Messaging**: Redpanda (optional event log backup).
- **Envoy Gateway Traffic Policies**: Envoy Gateway `Gateway` / `HTTPRoute` with WebSocket upgrade rules and long-lived client connection timeout tuning.

---

## Scenario 4: Crowdsourced Public Environmental Sensor Network

### Business Domain

A global, open-source community platform ingesting continuous environmental telemetry (air quality index, particulate matter $PM_{2.5}$, temperature) from distributed hardware monitors.

### Baseline Implementation

- **Architecture**: Monolithic CRUD application.
- **Behavior**: Ingests raw sensor values directly into a central database table, overwriting previous metrics or storing simple time stamps.

### Architecture Evolution

- **Phase 1: Hexagonal Architecture (Ports & Adapters)**
  - Decouple core domain logic (AQI calculation rules, safety threshold triggers) from external dependencies (databases, messaging brokers, HTTP frameworks).
  - Enable hot-swapping storage engines (e.g., switching PostgreSQL for InfluxDB) without altering core business rules.
- **Phase 2: Event Sourcing**
  - Store every sensor reading as an immutable, append-only Redpanda topic stream rather than mutating state.
  - Replay Redpanda event streams to reconstruct historical data states and provide verifiable data auditing.

### Kind & Envoy Infrastructure

- **Databases**: InfluxDB or TimescaleDB.
- **Messaging**: Redpanda or NATS.
- **Envoy Gateway Traffic Policies**: API routing rules using Kubernetes Gateway API (`HTTPRoute`) managed by Envoy Gateway to separate public contributor endpoints from internal processing services.

---

## Scenario 5: Decentralized Community Health Network

### Business Domain

A federated healthcare network connecting independent rural clinics, mobile health units, and regional hospitals to share anonymized diagnostic data and treatment outcomes for underserved populations.

### Baseline Implementation

- **Architecture**: Centralized API Gateway routing all requests to a single monolithic service.
- **Behavior**: All clinics connect directly to a central database, creating a single point of failure.

### Architecture Evolution

- **Phase 1: Load Balancing + Circuit Breaker Pattern**
  - Distribute traffic across multiple regional processing nodes using Envoy Gateway load balancing policies.
  - Implement circuit breakers to gracefully degrade when regional nodes become unavailable (protecting downstream services during network instability).
  - Add retry logic with exponential backoff for transient failures in satellite connections.
- **Phase 2: Scatter-Gather Pattern**
  - Query multiple independent clinic databases in parallel to aggregate patient treatment history.
  - Gather results with timeout handling for slow-responding remote clinics.
  - Return aggregated view even with partial data when some nodes are unreachable.

### Kind & Envoy Infrastructure

- **Databases**: PostgreSQL (per-region sharding).
- **Messaging**: Redpanda for async result aggregation.
- **Envoy Gateway Traffic Policies**: Envoy Gateway `BackendTrafficPolicy` for circuit breaking and connection pooling; `HTTPRoute` with weighted load balancing across regional services; timeout policies for scatter-gather request coordination.

---

## Scenario 6: Open Scientific Data Processing Platform

### Business Domain

A collaborative platform for processing large-scale climate research datasets (satellite imagery, ocean temperature sensors) contributed by universities and citizen scientists worldwide.

### Baseline Implementation

- **Architecture**: Sequential batch processing on a single server.
- **Behavior**: Datasets processed one at a time in a queue with no parallelization.

### Architecture Evolution

- **Phase 1: Map-Reduce Pattern**
  - Split large satellite image datasets across worker pods (Map phase).
  - Each worker applies data normalization and feature extraction in parallel.
  - Reduce phase aggregates processed chunks into unified climate models.
  - Use Redpanda to coordinate map tasks and collect reduce results.
- **Phase 2: Sidecar & Ambassador Pattern**
  - Deploy metric collection sidecars alongside each worker pod for transparent observability.
  - Use ambassador pattern to proxy external S3/object storage access with consistent auth.
  - Inject compression and caching logic via sidecars without modifying worker code.

### Kind & Envoy Infrastructure

- **Databases**: InfluxDB (time-series metrics) + S3-compatible MinIO (object storage).
- **Messaging**: Redpanda for task distribution and result collection.
- **Envoy Gateway Traffic Policies**: Envoy Gateway `Gateway` routes for external API access; sidecar proxies for telemetry collection.

---

## Scenario 7: Multi-Language Humanitarian Aid Portal

### Business Domain

A web portal serving displaced populations with aid resources, legal assistance, and community services, available in 20+ languages across mobile and desktop interfaces.

### Baseline Implementation

- **Architecture**: Single monolithic frontend serving all device types.
- **Behavior**: One-size-fits-all API responses with heavy client-side filtering and transformation.

### Architecture Evolution

- **Phase 1: Backends for Frontends (BFF) Pattern**
  - Create dedicated backend services optimized for mobile (low bandwidth), web (rich features), and SMS (text-only) interfaces.
  - Mobile BFF returns compressed JSON with minimal fields.
  - Web BFF includes full metadata and localization.
  - SMS BFF generates text-based responses.
- **Phase 2: Anti-Corruption Layer Pattern**
  - Integrate with legacy government assistance databases using anti-corruption adapters.
  - Translate between modern JSON APIs and legacy SOAP/XML services.
  - Protect core domain model from external system changes and versioning issues.

### Kind & Envoy Infrastructure

- **Databases**: PostgreSQL (modern domain model) + adapter layer for legacy systems.
- **Messaging**: Redpanda for async legacy system integration.
- **Envoy Gateway Traffic Policies**: Envoy Gateway `HTTPRoute` with path-based routing to different BFF services (`/mobile/*`, `/web/*`, `/sms/*`); header-based routing for API versioning.

---

## Scenario 8: Progressive Deployment for Emergency Alert System

### Business Domain

A critical national alert system for distributing emergency warnings (earthquakes, floods, wildfires) that must deploy updates without downtime while ensuring reliability.

### Baseline Implementation

- **Architecture**: Direct deployment with downtime during updates.
- **Behavior**: System goes offline during deployments, missing critical alert windows.

### Architecture Evolution

- **Phase 1: Blue-Green Deployment Pattern**
  - Maintain two identical production environments (Blue and Green).
  - Deploy updates to inactive environment first.
  - Switch traffic atomically using Envoy Gateway routing rules.
  - Instant rollback capability by switching back if issues detected.
- **Phase 2: Canary Deployment + Bulkhead Pattern**
  - Route 5% of traffic to new version first (canary).
  - Monitor error rates and latency before full rollout.
  - Implement bulkhead isolation to protect critical alert paths from non-critical features.
  - Use separate thread pools/connection pools for high-priority vs. low-priority alerts.

### Kind & Envoy Infrastructure

- **Databases**: PostgreSQL (alert history) + Redis (rate limiting state).
- **Messaging**: Redpanda for alert distribution pipeline.
- **Envoy Gateway Traffic Policies**: Envoy Gateway `HTTPRoute` with weighted traffic splitting for canary deployments; `ClientTrafficPolicy` for rate limiting and connection management; separate routing rules for blue/green environment switching.
