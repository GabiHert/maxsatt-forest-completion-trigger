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

## Scheduler Registration

This service receives events from the MaxSatt Scheduler Trigger via SNS/SQS. The scheduler uses DynamoDB TTL expiration to trigger events at 15-minute intervals.

### How It Works

1. A DynamoDB item with TTL is inserted into `maxsatt.scheduler-config` table
2. When TTL expires, DynamoDB deletes the item and triggers a stream event
3. The scheduler-trigger Lambda receives the stream event
4. Scheduler-trigger publishes to SNS topic `maxsatt-scheduler-trigger`
5. Forest-completion-trigger receives the message via its SQS subscription
6. After processing, scheduler-trigger creates a new DynamoDB item with the next TTL

### Initial Setup

To register the forest-completion-trigger with the scheduler, run:

```bash
# For development environment
./scripts/register-scheduler.sh dev

# For production environment
./scripts/register-scheduler.sh prd
```

### Prerequisites for Registration

- AWS CLI configured with credentials that have DynamoDB write access
- `jq` installed for JSON processing
- Access to the `maxsatt.scheduler-config` DynamoDB table

### Seed Data

The scheduler configuration is defined in `scripts/seed-scheduler-event.json`:

| Field | Value | Description |
|-------|-------|-------------|
| `event_id` | `forest-completion-check-001` | Unique identifier for this schedule |
| `service` | `MAXSATT_FOREST_COMPLETION_TRIGGER` | Service identifier |
| `event_type` | `FOREST_COMPLETION_CHECK` | Event type sent to the service |
| `schedule.execution_hours` | 96 time slots | Every 15 minutes (00:00, 00:15, ..., 23:45) |
| `schedule.interval_days` | 0 | Run every day |
| `schedule.ignore_week_days` | [] | No days skipped |
| `schedule.limit` | null | No end date (runs indefinitely) |

### Verifying Registration

After running the registration script, you can verify the schedule was created:

```bash
# Query DynamoDB for the schedule
aws dynamodb query \
  --table-name maxsatt.scheduler-config \
  --index-name event_id-index \
  --key-condition-expression "event_id = :eid" \
  --expression-attribute-values '{":eid": {"S": "forest-completion-check-001"}}' \
  --region us-east-1
```

### Canceling the Schedule

To cancel the schedule, you can either:

1. Delete the DynamoDB item directly
2. Set the `cancelled` field to `true`

```bash
# Find and delete the schedule
aws dynamodb query \
  --table-name maxsatt.scheduler-config \
  --index-name event_id-index \
  --key-condition-expression "event_id = :eid" \
  --expression-attribute-values '{":eid": {"S": "forest-completion-check-001"}}' \
  --region us-east-1 \
  --query 'Items[0].id.S' \
  --output text | xargs -I {} aws dynamodb delete-item \
    --table-name maxsatt.scheduler-config \
    --key '{"id": {"S": "{}"}}' \
    --region us-east-1
```

## Testing

- **BDD-First**: All features start with `.feature` files
- **Coverage**: 80% unit test coverage required
- **Mutation**: 84% mutation score required
