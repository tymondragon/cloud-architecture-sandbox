# Scenario 1: Natural Disaster Supply Matching System

## Pattern 1: Event-Driven Architecture (EDA)

### Phase 1: Conceptual Teaching

#### 1. Core Concept

**Event-Driven Architecture** is like a newsroom's breaking news alert system. When a reporter files a story (an event), they don't directly call every subscriber—instead, they publish it to a central wire service, and interested parties (editors, translators, fact-checkers) each pick it up independently and work in parallel. The reporter doesn't wait around or even know who's consuming their work; they just move on to the next story. In software, this means services publish events to a message broker and continue working, while other services consume those events asynchronously at their own pace.

#### 2. Problem vs. Solution

**Problem EDA Solves:**
- **Tight coupling**: In synchronous systems, Service A must wait for Service B, C, and D to finish before responding. If Service C is slow or down, everything stalls.
- **Cascading failures**: One slow database query blocks all upstream callers.
- **Scalability bottlenecks**: You can't independently scale services that are tightly synchronized.

**Solution & Trade-offs:**

✅ **Pros:**
- **Decoupling**: Publishers don't know (or care) who consumes their events
- **Scalability**: Add more consumers without touching producers
- **Resilience**: If a consumer is down, events queue up; system doesn't fail
- **Asynchronous processing**: API returns immediately; heavy work happens in background

⚠️ **Cons:**
- **Eventual consistency**: Data isn't immediately available everywhere (you give up real-time guarantees)
- **Complexity**: Debugging distributed async flows is harder than following synchronous call stacks
- **Message broker dependency**: Your broker (NATS, RabbitMQ, Redpanda) becomes a critical piece of infrastructure
- **Ordering challenges**: Events may arrive out of order without careful design

#### 3. Minimal Example (Go)

Here's a simplified Go example showing the pattern:

```go
package main

import (
    "encoding/json"
    "log"
    "github.com/nats-io/nats.go"
)

// Event: What happened in the system
type SupplyRequest struct {
    RequestID string `json:"request_id"`
    Location  string `json:"location"`
    Item      string `json:"item"`
    Quantity  int    `json:"quantity"`
}

// Publisher: API handler (returns immediately)
func HandleSupplyRequest(nc *nats.Conn, req SupplyRequest) error {
    // Serialize the event
    data, err := json.Marshal(req)
    if err != nil {
        return err
    }

    // Publish to message broker and return immediately
    // Notice: We DON'T wait for processing to complete!
    err = nc.Publish("supply.requested", data)
    if err != nil {
        return err
    }

    log.Printf("Published supply request %s - returning 202 Accepted", req.RequestID)
    return nil // API returns here; processing happens elsewhere
}

// Consumer: Worker pod listening to events
func SupplyAllocationWorker(nc *nats.Conn) {
    // Subscribe to events asynchronously
    nc.Subscribe("supply.requested", func(msg *nats.Msg) {
        var req SupplyRequest
        json.Unmarshal(msg.Data, &req)

        // Heavy processing happens here (database queries, ML algorithms, etc.)
        log.Printf("Worker processing request %s for %d x %s",
            req.RequestID, req.Quantity, req.Item)

        // This could take seconds or minutes - but API already returned!
        allocateSupplies(req)
    })
}

func allocateSupplies(req SupplyRequest) {
    // Complex allocation logic here
    // Database writes, inventory checks, routing algorithms, etc.
}
```

**Why this pattern?**
- `nc.Publish()` is **non-blocking** - the HTTP handler returns `202 Accepted` immediately
- The **worker** runs in a separate pod/process and can scale independently
- If the worker crashes, messages stay in the broker's queue (durability)
- You can add more workers without changing the publisher code

**Pointer syntax note** (since you're learning Go):
- `nc *nats.Conn` means `nc` is a **pointer** to a NATS connection object
- We use pointers here because connections are stateful resources that should be shared, not copied
- The `*` dereferences the pointer to access the actual object's methods

---

### Phase 2: Local Kind Cluster Blueprint

#### 1. Business Scenario Context

**Real-World Problem:**
During a natural disaster (earthquake, hurricane, flood), emergency relief organizations receive hundreds of supply requests from field agents, while donors simultaneously offer supplies. The system must:

- Accept requests even when infrastructure is unstable
- Match supplies to requests based on location, urgency, and availability
- Continue operating even if allocation services are temporarily down
- Handle traffic spikes without dropping critical requests

**Concrete Inputs & Outputs:**

**Input (HTTP POST):**
```json
{
  "request_id": "REQ-2024-001",
  "location": "shelter-downtown",
  "items": [
    {"type": "water_bottles", "quantity": 500},
    {"type": "medical_kits", "quantity": 20}
  ],
  "priority": "high"
}
```

**Immediate Output (202 Accepted):**
```json
{
  "request_id": "REQ-2024-001",
  "status": "accepted",
  "message": "Request queued for processing"
}
```

**Eventual Output (consumed from events, written to DB):**
- Allocation records showing which warehouse will fulfill the request
- Estimated delivery time
- Tracking information

---

#### 2. Infrastructure Layer (Helm-First)

##### Cluster Setup

```bash
# Create Kind cluster if not exists
kind create cluster --config cluster/kind-config.yaml

# Verify cluster
kubectl cluster-info --context kind-kind
kubectl get nodes
```

##### Add Helm Repositories

```bash
# Add required repositories
helm repo add envoy-gateway https://gateway.envoyproxy.io
helm repo add redpanda https://charts.redpanda.com/
helm repo add cnpg https://cloudnative-pg.io/charts/
helm repo update

# Verify repositories
helm repo list
```

##### Install Envoy Gateway (Mandatory)

```bash
helm install eg envoy-gateway/gateway-helm \
  --version v1.2.4 \
  --namespace envoy-gateway-system \
  --create-namespace \
  --wait
```

**Why Envoy Gateway?**
- Implements Kubernetes Gateway API (modern standard)
- Provides rate limiting to prevent request drops during traffic spikes
- Supports advanced traffic policies for resilience

Wait for Envoy Gateway to be ready:
```bash
kubectl wait --timeout=5m -n envoy-gateway-system \
  deployment/envoy-gateway \
  --for=condition=Available
```

##### Install Redpanda (Event Broker)

**Why Redpanda for this scenario?**
- Kafka-compatible API (industry standard for event streaming)
- Designed for high-throughput event logs
- Better performance than traditional Kafka (no JVM overhead)
- Perfect for saga patterns (Phase 2) which need event sourcing capabilities

Create Helm values file at `infra/redpanda/values.yaml`:

```yaml
# Redpanda configuration for local Kind cluster
auth:
  sasl:
    enabled: false  # Simplified for learning; enable in production

storage:
  persistentVolume:
    enabled: true
    size: 10Gi
    storageClass: standard

resources:
  cpu:
    cores: 1
  memory:
    container:
      max: 2Gi

# Single-node setup for local development
statefulset:
  replicas: 1

# Enable external connectivity
external:
  enabled: false  # We'll use port-forward for local access

# Console for web UI
console:
  enabled: true
```

Install Redpanda:

```bash
helm install redpanda redpanda/redpanda \
  --namespace redpanda \
  --create-namespace \
  -f infra/redpanda/values.yaml \
  --wait
```

Verify installation:
```bash
kubectl get pods -n redpanda -w
# Wait until STATUS is Running
```

##### Install CloudNativePG Operator (PostgreSQL)

**Why CloudNativePG?**
- CNCF Sandbox project (cloud-native best practices)
- Uses Kubernetes operators (declarative cluster management)
- Better than traditional Helm charts for stateful workloads
- Automatic failover, backup, and recovery capabilities

Install the CloudNativePG operator:

```bash
helm install cnpg cnpg/cloudnative-pg \
  --namespace cnpg-system \
  --create-namespace \
  --wait
```

Verify operator is running:
```bash
kubectl get pods -n cnpg-system
```

Create PostgreSQL Cluster manifest at `infra/postgresql/cluster.yaml`:

```yaml
# PostgreSQL Cluster using CloudNativePG operator
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: supply-db
  namespace: databases
spec:
  instances: 1  # Single instance for local development

  # Database initialization
  bootstrap:
    initdb:
      database: supply_matching
      owner: supply_user
      secret:
        name: supply-db-credentials

  # Storage configuration
  storage:
    size: 5Gi
    storageClass: standard

  # Resource limits for local development
  resources:
    requests:
      memory: "256Mi"
      cpu: "250m"
    limits:
      memory: "512Mi"
      cpu: "500m"

  # Monitoring
  monitoring:
    enablePodMonitor: false
---
# Secret for database credentials
apiVersion: v1
kind: Secret
metadata:
  name: supply-db-credentials
  namespace: databases
type: kubernetes.io/basic-auth
stringData:
  username: supply_user
  password: local-dev-password  # Use sealed secrets in production!
---
# Namespace
apiVersion: v1
kind: Namespace
metadata:
  name: databases
```

Deploy the PostgreSQL cluster:

```bash
kubectl apply -f infra/postgresql/cluster.yaml
```

Wait for cluster to be ready:
```bash
kubectl wait --for=condition=Ready cluster/supply-db -n databases --timeout=5m
```

---

#### 3. Incremental Layering

**What we're building on:**
- This is the **first scenario**, so we're establishing the baseline infrastructure
- Future scenarios will reuse Envoy Gateway, and potentially Redpanda/PostgreSQL

**What's new in this scenario:**
- Redpanda broker (shared across future event-driven scenarios)
- PostgreSQL database (shared across scenarios needing relational storage)
- Two microservices: API Gateway (publisher) and Allocation Worker (consumer)

**What can be reused later:**
- Scenario 3 (Food Redistribution) can reuse PostgreSQL for writes
- Scenario 4 (Environmental Sensors) can reuse Redpanda for event sourcing

---

#### 4. Component Breakdown

We'll build **two microservices** in Go:

##### Service 1: `supply-api` (Event Publisher)

**Responsibility:**
- Receive HTTP POST requests with supply needs
- Validate input (basic schema validation)
- Publish event to Redpanda topic `supply.requests`
- Return `202 Accepted` immediately

**Core Business Logic:**
```go
// Accept request, publish event, return immediately
func (h *Handler) CreateSupplyRequest(w http.ResponseWriter, r *http.Request) {
    var req SupplyRequest

    // Parse and validate
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Publish to Redpanda (non-blocking)
    if err := h.publisher.PublishSupplyRequest(ctx, req); err != nil {
        http.Error(w, "Failed to queue request", http.StatusServiceUnavailable)
        return
    }

    // Return immediately - processing happens asynchronously
    w.WriteHeader(http.StatusAccepted)
    json.NewEncoder(w).Encode(map[string]string{
        "request_id": req.RequestID,
        "status": "accepted",
    })
}
```

**Why Go?**
- Excellent Kafka/Redpanda client libraries (`franz-go`)
- Fast HTTP performance with `net/http`
- Simple concurrency model for handling high traffic

**Key Interface:**
```go
type EventPublisher interface {
    PublishSupplyRequest(ctx context.Context, req SupplyRequest) error
}
```

##### Service 2: `allocation-worker` (Event Consumer)

**Responsibility:**
- Subscribe to `supply.requests` topic in Redpanda
- Process each event: match supplies with requests using allocation algorithm
- Write allocation results to PostgreSQL
- Acknowledge message after successful processing

**Core Business Logic:**
```go
func (w *Worker) ProcessSupplyRequest(ctx context.Context, event SupplyRequestEvent) error {
    // Complex allocation logic
    allocation := w.allocator.FindBestMatch(event.Location, event.Items)

    // Write to database (transactional)
    if err := w.db.SaveAllocation(ctx, allocation); err != nil {
        return err // Message will be retried
    }

    // Only ack if database write succeeded
    return nil
}
```

**Key Interface:**
```go
type SupplyAllocator interface {
    FindBestMatch(location string, items []Item) (*Allocation, error)
}
```

---

#### 5. Gateway API Resources

##### Gateway Resource

Create `infra/envoy-gateway/gateway.yaml`:

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: main-gateway
  namespace: envoy-gateway-system
spec:
  gatewayClassName: eg
  listeners:
    - name: http
      protocol: HTTP
      port: 80
      allowedRoutes:
        namespaces:
          from: All
```

Apply the Gateway:
```bash
kubectl apply -f infra/envoy-gateway/gateway.yaml
```

##### HTTPRoute for Supply API

Create `manifests/httproute.yaml`:

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: supply-api-route
  namespace: scenario-01
spec:
  parentRefs:
    - name: main-gateway
      namespace: envoy-gateway-system
  hostnames:
    - "supply-api.local"
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: /api/v1/supplies
      backendRefs:
        - name: supply-api
          port: 8080
```

##### Rate Limiting Policy (Critical for Emergency Scenarios)

Create `manifests/rate-limit-policy.yaml`:

```yaml
apiVersion: gateway.envoyproxy.io/v1alpha1
kind: ClientTrafficPolicy
metadata:
  name: supply-api-rate-limit
  namespace: envoy-gateway-system
spec:
  targetRef:
    group: gateway.networking.k8s.io
    kind: Gateway
    name: main-gateway
  rateLimit:
    type: Global
    global:
      rules:
        - limit:
            requests: 100
            unit: Second
```

**Why rate limiting?**
During disasters, traffic can spike unpredictably. Instead of overwhelming the system and dropping requests, rate limiting:
- Queues excess requests in Envoy
- Prevents database connection pool exhaustion
- Returns `429 Too Many Requests` instead of silent failures

---

#### 6. Application Manifests

##### Supply API Deployment

Create `manifests/supply-api-deployment.yaml`:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: scenario-01
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: supply-api
  namespace: scenario-01
spec:
  replicas: 2
  selector:
    matchLabels:
      app: supply-api
  template:
    metadata:
      labels:
        app: supply-api
    spec:
      containers:
        - name: api
          image: supply-api:latest
          imagePullPolicy: Never  # Use local image in Kind
          ports:
            - containerPort: 8080
          env:
            - name: REDPANDA_BROKERS
              value: "redpanda-0.redpanda.redpanda.svc.cluster.local:9092"
            - name: TOPIC_NAME
              value: "supply.requests"
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "256Mi"
              cpu: "200m"
---
apiVersion: v1
kind: Service
metadata:
  name: supply-api
  namespace: scenario-01
spec:
  selector:
    app: supply-api
  ports:
    - protocol: TCP
      port: 8080
      targetPort: 8080
```

##### Allocation Worker Deployment

Create `manifests/allocation-worker-deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: allocation-worker
  namespace: scenario-01
spec:
  replicas: 3  # Scale consumers independently!
  selector:
    matchLabels:
      app: allocation-worker
  template:
    metadata:
      labels:
        app: allocation-worker
    spec:
      containers:
        - name: worker
          image: allocation-worker:latest
          imagePullPolicy: Never
          env:
            - name: REDPANDA_BROKERS
              value: "redpanda-0.redpanda.redpanda.svc.cluster.local:9092"
            - name: TOPIC_NAME
              value: "supply.requests"
            - name: CONSUMER_GROUP
              value: "allocation-workers"
            - name: POSTGRES_HOST
              value: "supply-db-rw.databases.svc.cluster.local"  # CloudNativePG read-write service
            - name: POSTGRES_DB
              value: "supply_matching"
            - name: POSTGRES_USER
              valueFrom:
                secretKeyRef:
                  name: supply-db-credentials
                  key: username
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: supply-db-credentials
                  key: password
          resources:
            requests:
              memory: "256Mi"
              cpu: "200m"
            limits:
              memory: "512Mi"
              cpu: "500m"
```

**Notice:** We have 3 worker replicas but 2 API replicas. This is intentional—we can scale consumers based on processing needs without touching the API layer. This is the power of decoupling!

---

#### 7. Observability & Connectivity

##### Access the Supply API

**Option 1: Port-forward (Quick Testing)**

```bash
# Forward Envoy Gateway service
kubectl port-forward -n envoy-gateway-system \
  service/envoy-main-gateway 8080:80
```

Then test with:
```bash
curl -X POST http://localhost:8080/api/v1/supplies \
  -H "Host: supply-api.local" \
  -H "Content-Type: application/json" \
  -d '{
    "request_id": "REQ-001",
    "location": "shelter-downtown",
    "items": [
      {"type": "water_bottles", "quantity": 500}
    ],
    "priority": "high"
  }'
```

**Option 2: Update /etc/hosts (More Realistic)**

```bash
# Get the port-forward address
echo "127.0.0.1 supply-api.local" | sudo tee -a /etc/hosts
```

Then:
```bash
curl -X POST http://supply-api.local:8080/api/v1/supplies \
  -H "Content-Type: application/json" \
  -d '{"request_id": "REQ-001", ...}'
```

##### Monitor Redpanda Topics

```bash
# List topics
kubectl exec -n redpanda redpanda-0 -- \
  rpk topic list

# Create topic manually (if not auto-created)
kubectl exec -n redpanda redpanda-0 -- \
  rpk topic create supply.requests --partitions 3 --replicas 1

# Consume messages (see what workers are processing)
kubectl exec -n redpanda redpanda-0 -- \
  rpk topic consume supply.requests --format='%v\n'

# Check consumer group lag (are workers keeping up?)
kubectl exec -n redpanda redpanda-0 -- \
  rpk group describe allocation-workers
```

##### View Application Logs

```bash
# API logs (publishers)
kubectl logs -n scenario-01 -l app=supply-api -f --tail=50

# Worker logs (consumers)
kubectl logs -n scenario-01 -l app=allocation-worker -f --tail=50

# Specific pod
kubectl logs -n scenario-01 allocation-worker-5d7b8c9f4-abc12 -f
```

##### Check PostgreSQL Data

```bash
# Connect to PostgreSQL (CloudNativePG cluster)
kubectl exec -it -n databases supply-db-1 -- \
  psql -U supply_user -d supply_matching

# Query allocations
SELECT * FROM allocations ORDER BY created_at DESC LIMIT 10;

# Exit with \q
```

---

#### 8. Validation Steps

Let's verify the Event-Driven Architecture is working end-to-end.

##### Step 1: Verify Infrastructure

```bash
# Check all pods are running
kubectl get pods -n redpanda
kubectl get pods -n databases
kubectl get pods -n cnpg-system
kubectl get pods -n envoy-gateway-system
kubectl get pods -n scenario-01

# Should see:
# - redpanda-0 (Running)
# - supply-db-1 (Running) - CloudNativePG cluster
# - cnpg-cloudnative-pg-* (Running) - Operator
# - envoy-gateway-* (Running)
# - supply-api-* (2 replicas Running)
# - allocation-worker-* (3 replicas Running)
```

##### Step 2: Create Redpanda Topic

```bash
kubectl exec -n redpanda redpanda-0 -- \
  rpk topic create supply.requests --partitions 3 --replicas 1
```

##### Step 3: Send Test Request

```bash
# Port-forward Envoy Gateway
kubectl port-forward -n envoy-gateway-system \
  service/envoy-main-gateway 8080:80 &

# Send supply request
curl -X POST http://localhost:8080/api/v1/supplies \
  -H "Host: supply-api.local" \
  -H "Content-Type: application/json" \
  -d '{
    "request_id": "TEST-001",
    "location": "shelter-downtown",
    "items": [
      {"type": "water_bottles", "quantity": 500},
      {"type": "blankets", "quantity": 100}
    ],
    "priority": "high"
  }'

# Expected response (IMMEDIATE 202):
# {"request_id":"TEST-001","status":"accepted"}
```

##### Step 4: Verify Event in Redpanda

```bash
# Consume from topic (should see the event)
kubectl exec -n redpanda redpanda-0 -- \
  rpk topic consume supply.requests --num 1

# Output should show:
# {
#   "request_id": "TEST-001",
#   "location": "shelter-downtown",
#   ...
# }
```

##### Step 5: Check Worker Processing

```bash
# Watch worker logs
kubectl logs -n scenario-01 -l app=allocation-worker --tail=20

# Should see:
# "Processing supply request TEST-001"
# "Allocated supplies from warehouse-north"
# "Saved allocation to database"
```

##### Step 6: Verify Database Write

```bash
# Query PostgreSQL (CloudNativePG cluster)
kubectl exec -it -n databases supply-db-1 -- \
  psql -U supply_user -d supply_matching -c \
  "SELECT request_id, location, status, allocated_warehouse FROM allocations WHERE request_id='TEST-001';"

# Should show:
#  request_id | location          | status    | allocated_warehouse
# ------------+-------------------+-----------+--------------------
#  TEST-001   | shelter-downtown  | allocated | warehouse-north
```

##### Step 7: Demonstrate Async Behavior

Send multiple requests rapidly:

```bash
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/v1/supplies \
    -H "Host: supply-api.local" \
    -H "Content-Type: application/json" \
    -d "{\"request_id\":\"BULK-$i\",\"location\":\"shelter-$i\",\"items\":[{\"type\":\"water\",\"quantity\":100}],\"priority\":\"medium\"}" &
done
wait

# All requests return 202 immediately
# Workers process them asynchronously (check logs)
kubectl logs -n scenario-01 -l app=allocation-worker --tail=50
```

##### Step 8: Demonstrate Consumer Group Rebalancing

```bash
# Scale workers up
kubectl scale deployment -n scenario-01 allocation-worker --replicas=5

# Check consumer group rebalancing
kubectl exec -n redpanda redpanda-0 -- \
  rpk group describe allocation-workers

# You'll see partitions redistributed across 5 consumers
```

---

#### 9. Clean Teardown & Isolation

##### Remove Scenario Workloads (Keep Infrastructure)

```bash
# Delete application deployments
kubectl delete namespace scenario-01

# This removes:
# - supply-api deployment & service
# - allocation-worker deployment
# - HTTPRoute

# Redpanda topic cleanup (optional - keeps data for replay)
kubectl exec -n redpanda redpanda-0 -- \
  rpk topic delete supply.requests
```

##### Preserve Shared Infrastructure

**DO NOT DELETE** (reusable for future scenarios):
```bash
# Keep these running:
# - Envoy Gateway (envoy-gateway-system namespace)
# - CloudNativePG Operator (cnpg-system namespace)
# - Redpanda (redpanda namespace)
# - PostgreSQL Cluster (databases namespace)
```

##### Full Teardown (Between Learning Sessions)

Only if starting completely fresh:

```bash
# Remove everything
kubectl delete -f infra/postgresql/cluster.yaml  # Delete PostgreSQL cluster
helm uninstall redpanda -n redpanda
helm uninstall cnpg -n cnpg-system
helm uninstall eg -n envoy-gateway-system

# Delete namespaces
kubectl delete namespace redpanda databases cnpg-system scenario-01

# Or delete entire cluster
kind delete cluster
```

---

## Summary: What We Built

✅ **Event-Driven Architecture Implementation:**
- API accepts requests and publishes events (decoupled)
- Workers consume events asynchronously
- System returns `202 Accepted` immediately (no blocking)

✅ **Infrastructure:**
- Envoy Gateway with rate limiting
- Redpanda for event streaming
- PostgreSQL for state storage
- All deployed with Helm

✅ **Key Learnings:**
- Publishers don't wait for consumers
- Consumers can scale independently
- Message acknowledgment prevents data loss
- Eventual consistency trade-off

---

## Next Pattern: Saga Pattern

The Saga Pattern builds on this EDA foundation to handle **distributed transactions** across multiple services (e.g., Reserve Supplies → Allocate Transport → Update Shelter Queue). We'll add compensating transactions for rollback when steps fail.
