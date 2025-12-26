#language: en
#utf-8

@all @forest_completion
Feature: Forest Completion Detection and Notification

  Background:
    Given the "LOG_LEVEL" env var is set to "debug"
    And the "SERVICE_NAME" env var is set to "maxsatt-forest-completion-trigger"
    And the "FOREST_EVENTS_TOPIC_ARN" env var is set to "arn:aws:sns:us-east-1:123456789:forest-events-topic"
    And the "MAXSATT_API_URL" env var is set to the mocked api url
    And the "MAXSATT_AUTH_URL" env var is set to the mocked oauth url
    And the "MAXSATT_CLIENT_ID" env var is set to "test-client-id"
    And the "MAXSATT_CLIENT_SECRET" env var is set to "test-client-secret"
    And the "MAX_FORESTS_PER_INVOCATION" env var is set to "100"
    And the "DISCORD_WEBHOOK_URL" env var is set to the mocked api url
    And the oauth server returns a valid token

  # ========================================
  # SUCCESS SCENARIOS
  # ========================================

  @success @sns_integration
  Scenario: Forest with all processings completed triggers notification
    Given the -1 "GET" request to "/v1/processings" returns status 200 with the following response
    """
    {
      "data": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440001",
          "field_id": "5a801269-110d-46c0-89dd-09a6dc195401",
          "status": "COMPLETED",
          "notified_at": null,
          "field": {
            "id": "5a801269-110d-46c0-89dd-09a6dc195401",
            "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c0001",
            "name": "Field A"
          }
        },
        {
          "id": "550e8400-e29b-41d4-a716-446655440002",
          "field_id": "5a801269-110d-46c0-89dd-09a6dc195402",
          "status": "COMPLETED",
          "notified_at": null,
          "field": {
            "id": "5a801269-110d-46c0-89dd-09a6dc195402",
            "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c0001",
            "name": "Field B"
          }
        }
      ]
    }
    """
    And the -1 "PUT" request to "/v1/processing/*" returns status 200 with the following response
    """
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "status": "COMPLETED",
      "notified_at": "2025-12-20T10:00:00Z"
    }
    """
    When the following event is received via sqs
    """
    {
      "trigger": "scheduled"
    }
    """
    Then the lambda should finish without errors
    And the sns topic "arn:aws:sns:us-east-1:123456789:forest-events-topic" should have 1 messages published
    And the sns topic "arn:aws:sns:us-east-1:123456789:forest-events-topic" should have a message published with "event_type" field equal to "NOTIFY"
    And the sns topic "arn:aws:sns:us-east-1:123456789:forest-events-topic" should have a message published with "event_data.forest_id" field equal to "87428a25-f0d2-4dba-acfd-e3d7320c0001"
    And the sns topic "arn:aws:sns:us-east-1:123456789:forest-events-topic" should have a message published with "source" field equal to "MAXSATT_FOREST_COMPLETION_TRIGGER"
    And the 0 "PUT" request for "/v1/processing/*" json body should have the variable "notified_at" with value "not nil"

  @success @multiple_forests
  Scenario: Multiple forests with completed processings trigger multiple notifications
    Given the -1 "GET" request to "/v1/processings" returns status 200 with the following response
    """
    {
      "data": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440001",
          "field_id": "5a801269-110d-46c0-89dd-09a6dc195401",
          "status": "COMPLETED",
          "notified_at": null,
          "field": {
            "id": "5a801269-110d-46c0-89dd-09a6dc195401",
            "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c0001",
            "name": "Field A"
          }
        },
        {
          "id": "550e8400-e29b-41d4-a716-446655440002",
          "field_id": "5a801269-110d-46c0-89dd-09a6dc195402",
          "status": "COMPLETED",
          "notified_at": null,
          "field": {
            "id": "5a801269-110d-46c0-89dd-09a6dc195402",
            "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c0002",
            "name": "Field B"
          }
        }
      ]
    }
    """
    And the -1 "PUT" request to "/v1/processing/*" returns status 200 with the following response
    """
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "status": "COMPLETED",
      "notified_at": "2025-12-20T10:00:00Z"
    }
    """
    When the following event is received via sqs
    """
    {
      "trigger": "scheduled"
    }
    """
    Then the lambda should finish without errors
    And the sns topic "arn:aws:sns:us-east-1:123456789:forest-events-topic" should have 2 messages published

  # ========================================
  # SKIP SCENARIOS - No action needed
  # ========================================

  @skip @no_completed_forests
  Scenario: No forests with all processings completed returns empty result
    Given the -1 "GET" request to "/v1/processings" returns status 200 with the following response
    """
    {
      "data": []
    }
    """
    When the following event is received via sqs
    """
    {
      "trigger": "scheduled"
    }
    """
    Then the lambda should finish without errors
    And the sns topic "arn:aws:sns:us-east-1:123456789:forest-events-topic" should have 0 messages published

  @skip @already_notified
  Scenario: No processings ready for notification returns empty result
    Given the -1 "GET" request to "/v1/processings" returns status 200 with the following response
    """
    {
      "data": []
    }
    """
    When the following event is received via sqs
    """
    {
      "trigger": "scheduled"
    }
    """
    Then the lambda should finish without errors
    And the sns topic "arn:aws:sns:us-east-1:123456789:forest-events-topic" should have 0 messages published

  # ========================================
  # ERROR SCENARIOS
  # ========================================

  @error @api_error
  Scenario: API returns error when fetching processings
    Given the -1 "GET" request to "/v1/processings" returns status 500 with the following response
    """
    {
      "error": {
        "code": "INTERNAL_ERROR",
        "description": "Database connection failed"
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "trigger": "scheduled"
    }
    """
    Then the lambda should finish with "internal server" error

  @error @auth_error
  Scenario: API returns unauthorized when token is invalid
    Given the oauth server returns an invalid token response
    When the following event is received via sqs
    """
    {
      "trigger": "scheduled"
    }
    """
    Then the lambda should finish with "Invalid credentials" error
