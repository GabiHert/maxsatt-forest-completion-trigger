# MaxSatt Forest Completion Trigger

Lambda service that detects when all processings for a forest are completed and publishes notification events.

## Overview

This service runs every 15 minutes via the scheduler-trigger and:

1. Queries PostgreSQL to find forests where ALL processings across ALL fields have status = COMPLETED and at least one has `notified_at IS NULL`
2. Publishes to `forest-events-topic` with `event_type: "notify"`
3. Updates `notified_at` timestamp for all related processings
4. Limits to 100 forests per invocation
5. Uses DLQ + Discord notification on failure

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Scheduler Trigger                             │
│                    (every 15 min)                                │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│               Forest Completion Trigger                          │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────────┐   │
│  │ Lambda       │───▶│ PostgreSQL   │───▶│ SNS Publisher    │   │
│  │ Handler      │    │ Repository   │    │ (forest-events)  │   │
│  └──────────────┘    └──────────────┘    └──────────────────┘   │
│                              │                    │              │
│                              ▼                    ▼              │
│                    ┌──────────────┐     ┌──────────────────┐    │
│                    │ Update       │     │ DLQ + Discord    │    │
│                    │ notified_at  │     │ (on failure)     │    │
│                    └──────────────┘     └──────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
```

## Message Format

```json
{
  "spec_version": 1,
  "event_id": "uuid",
  "event_correlation_id": "uuid",
  "event_date": "2025-12-23T10:00:00Z",
  "source": "MAXSATT_FOREST_COMPLETION_TRIGGER",
  "event_type": "notify",
  "event_data": {
    "forest_id": "uuid",
    "processing_ids": ["uuid1", "uuid2"]
  }
}
```

## Development

### Prerequisites

- Go 1.25+
- Docker (for local testing)
- AWS CLI configured
- PostgreSQL database

### Commands

```bash
make run           # Run locally
make build-binary  # Build Lambda artifact
make validate      # Full validation (tests, mutation, deps)
make tests         # Unit tests
make test-bdd      # BDD tests (Godog)
make mutation-test # Mutation testing
make sort-imports  # Format imports
```

## Clean Architecture

```
internal/
├── domain/           # Business entities and errors
├── application/      # Use cases and adapter interfaces
├── integration/      # Controllers, repositories, publishers
└── infra/            # Framework setup, DI, server config
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `DB_HOST` | PostgreSQL host |
| `DB_PORT` | PostgreSQL port |
| `DB_NAME` | Database name |
| `DB_USERNAME` | Database username |
| `DB_PASSWORD` | Database password |
| `AWS_SNS_FOREST_EVENTS_TOPIC_ARN` | SNS topic for forest events |
| `AWS_SQS_DLQ_URL` | Dead letter queue URL |
| `DISCORD_WEBHOOK_URL` | Discord webhook for error notifications |
| `SERVICE_NAME` | Service identifier |

## Testing

- **BDD-First**: All features start with `.feature` files
- **Coverage**: 80% unit test coverage required
- **Mutation**: 84% mutation score required
