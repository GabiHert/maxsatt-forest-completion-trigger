#!/bin/bash
#
# Register Forest Completion Trigger with MaxSatt Scheduler
#
# This script inserts the scheduler configuration into DynamoDB
# to enable 15-minute interval scheduling for the forest completion check.
#
# Usage:
#   ./register-scheduler.sh <environment>
#
# Arguments:
#   environment - Target environment: dev or prd
#
# Requirements:
#   - AWS CLI configured with appropriate credentials
#   - Access to the maxsatt.scheduler-config DynamoDB table
#
# Examples:
#   ./register-scheduler.sh dev
#   ./register-scheduler.sh prd

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SEED_FILE="${SCRIPT_DIR}/seed-scheduler-event.json"

# Configuration
DYNAMODB_TABLE_NAME="maxsatt.scheduler-config"
AWS_REGION="${AWS_REGION:-us-east-1}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

usage() {
    echo "Usage: $0 <environment>"
    echo ""
    echo "Arguments:"
    echo "  environment    Target environment (dev or prd)"
    echo ""
    echo "Examples:"
    echo "  $0 dev"
    echo "  $0 prd"
    exit 1
}

calculate_next_ttl() {
    # Calculate the next 15-minute interval in Unix timestamp
    # TTL is when DynamoDB will delete the item, triggering the stream event
    local current_timestamp
    local current_minute
    local current_second
    local next_interval_minute
    local seconds_to_add
    local next_ttl

    current_timestamp=$(date +%s)
    current_minute=$(date +%M | sed 's/^0//')
    current_second=$(date +%S | sed 's/^0//')

    # Find next 15-minute mark (00, 15, 30, 45)
    case $(( current_minute % 15 )) in
        0)
            if [ "${current_second:-0}" -gt 0 ]; then
                seconds_to_add=$(( (15 - (current_minute % 15)) * 60 - current_second ))
            else
                seconds_to_add=$(( 15 * 60 ))
            fi
            ;;
        *)
            seconds_to_add=$(( (15 - (current_minute % 15)) * 60 - current_second ))
            ;;
    esac

    # Add a small buffer (30 seconds) to ensure we don't miss the window
    next_ttl=$(( current_timestamp + seconds_to_add + 30 ))

    echo "${next_ttl}"
}

check_existing_schedule() {
    local env="$1"
    local event_id="forest-completion-check-001"

    log_info "Checking for existing schedule..."

    local result
    result=$(aws dynamodb query \
        --table-name "${DYNAMODB_TABLE_NAME}" \
        --index-name "event_id-index" \
        --key-condition-expression "event_id = :eid" \
        --expression-attribute-values "{\":eid\": {\"S\": \"${event_id}\"}}" \
        --region "${AWS_REGION}" \
        --output json 2>/dev/null || echo '{"Count": 0}')

    local count
    count=$(echo "${result}" | jq -r '.Count // 0')

    if [ "${count}" -gt 0 ]; then
        return 0  # Schedule exists
    fi
    return 1  # Schedule does not exist
}

delete_existing_schedule() {
    local event_id="forest-completion-check-001"

    log_info "Querying existing schedules to delete..."

    local result
    result=$(aws dynamodb query \
        --table-name "${DYNAMODB_TABLE_NAME}" \
        --index-name "event_id-index" \
        --key-condition-expression "event_id = :eid" \
        --expression-attribute-values "{\":eid\": {\"S\": \"${event_id}\"}}" \
        --region "${AWS_REGION}" \
        --output json 2>/dev/null)

    local ids
    ids=$(echo "${result}" | jq -r '.Items[].id.S // empty')

    for id in ${ids}; do
        log_info "Deleting existing schedule with id: ${id}"
        aws dynamodb delete-item \
            --table-name "${DYNAMODB_TABLE_NAME}" \
            --key "{\"id\": {\"S\": \"${id}\"}}" \
            --region "${AWS_REGION}"
    done
}

register_schedule() {
    local env="$1"
    local ttl="$2"

    log_info "Reading seed file from: ${SEED_FILE}"

    if [ ! -f "${SEED_FILE}" ]; then
        log_error "Seed file not found: ${SEED_FILE}"
        exit 1
    fi

    # Generate a unique ID for this schedule item
    local unique_id
    unique_id="forest-completion-check-$(date +%Y%m%d%H%M%S)"

    # Create the item JSON with calculated TTL and unique ID
    local item_json
    item_json=$(cat "${SEED_FILE}" | \
        sed "s/PLACEHOLDER_TTL/${ttl}/g" | \
        jq --arg id "${unique_id}" '.id.S = $id')

    log_info "Registering schedule with DynamoDB..."
    log_info "Table: ${DYNAMODB_TABLE_NAME}"
    log_info "Region: ${AWS_REGION}"
    log_info "Environment: ${env}"
    log_info "Item ID: ${unique_id}"
    log_info "TTL: ${ttl} ($(date -r ${ttl} 2>/dev/null || date -d @${ttl} 2>/dev/null || echo 'unable to format'))"

    # Insert the item into DynamoDB
    aws dynamodb put-item \
        --table-name "${DYNAMODB_TABLE_NAME}" \
        --item "${item_json}" \
        --region "${AWS_REGION}"

    if [ $? -eq 0 ]; then
        log_info "Schedule registered successfully!"
        echo ""
        log_info "The scheduler will trigger at the next 15-minute interval."
        log_info "After TTL expiration, the DynamoDB stream will trigger the scheduler-trigger Lambda."
        log_info "The scheduler-trigger will then publish to SNS, and forest-completion-trigger will receive via SQS."
    else
        log_error "Failed to register schedule"
        exit 1
    fi
}

main() {
    if [ $# -lt 1 ]; then
        usage
    fi

    local env="$1"

    # Validate environment
    case "${env}" in
        dev|prd)
            ;;
        *)
            log_error "Invalid environment: ${env}. Must be 'dev' or 'prd'."
            usage
            ;;
    esac

    log_info "Starting Forest Completion Trigger scheduler registration..."
    log_info "Environment: ${env}"

    # Check AWS CLI is available
    if ! command -v aws &> /dev/null; then
        log_error "AWS CLI is not installed or not in PATH"
        exit 1
    fi

    # Check jq is available
    if ! command -v jq &> /dev/null; then
        log_error "jq is not installed or not in PATH"
        exit 1
    fi

    # Calculate next TTL
    local next_ttl
    next_ttl=$(calculate_next_ttl)
    log_info "Calculated next TTL: ${next_ttl}"

    # Check for existing schedule
    if check_existing_schedule "${env}"; then
        log_warn "Existing schedule found for forest-completion-check"
        read -p "Do you want to delete existing schedule and create a new one? (y/N): " confirm
        if [[ "${confirm}" =~ ^[Yy]$ ]]; then
            delete_existing_schedule
        else
            log_info "Keeping existing schedule. Exiting."
            exit 0
        fi
    fi

    # Register the schedule
    register_schedule "${env}" "${next_ttl}"
}

main "$@"
