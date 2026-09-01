# Claude Code Skills

This directory contains custom skills for the cloud-architecture-sandbox repository.

## Available Skills

### `/learn-pattern`

Start an interactive architecture pattern learning session with structured two-phase teaching.

**Usage:**
```bash
# List all available scenarios
/learn-pattern

# Start scenario 1 (Natural Disaster Supply Matching - EDA → Saga)
/learn-pattern 1

# Start scenario 2 (Wildlife Detection - Pipes & Filters → Transactional Outbox)
/learn-pattern 2

# ... and so on for scenarios 3-8
```

**What it does:**
- **Phase 1**: Teaches the pattern conceptually with analogies, trade-offs, and code examples
- **Phase 2**: Guides you through deploying a working implementation in Kind with Helm
- Provides exact commands for infrastructure setup, validation, and teardown
- Tracks progress through multi-pattern scenarios
- Follows best practices from CLAUDE.md and architecture-context.md

**Features:**
- Helm-first infrastructure deployment
- Appropriate messaging selection (NATS, RabbitMQ, Redpanda, Redis Pub/Sub)
- Envoy Gateway for all ingress/routing
- Incremental learning that builds on previous scenarios
- Go-friendly examples with explanations
- OTEL observability best practices

See `skills/learn-pattern.md` for full documentation.
