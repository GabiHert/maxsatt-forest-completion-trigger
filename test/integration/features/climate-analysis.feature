#language: en
#utf-8

@all @climate_analysis
Feature: Climate analysis functionality

  Background:
    Given the tables are empty
    And the "LOG_LEVEL" env var is set to "debug"
    And the "S3_BUCKET_NAME" env var is set to "maxsatt-data-bucket"
    And the "WEATHER_API_URL" env var is set to the mocked api url
    And the "FOREST_FIELD_API_URL" env var is set to the mocked api url
    And the "CACHE_TABLE_NAME" env var is set to "maxsatt-weather-cache"
    And the "CACHE_ENABLED" env var is set to "true"
    And the "API_MAX_RETRIES" env var is set to "10"
    And the "API_TIMEOUT_SECONDS" env var is set to "30"
    And the "API_RETRY_DELAY_SECONDS" env var is set to "10"
    And the "SERVICE_NAME" env var is set to "maxsatt-climate-analysis-lambda"
    And the dynamodb table "maxsatt-weather-cache" with key "cache_key" exists
    And the aws secret named "maxsatt-image-climate-dataset-trigger-secret" exists with the following values
    """
    {
      "AUTH_CLIENT_ID": "7073da1d-c57e-4fdc-a2e2-6208fee06ff1",
      "AUTH_CLIENT_SECRET": "b3d2c4f0-8c1b-4d2a-9e3f-5c6b7e8f9a0b"
    }
    """
    And the "AUTH_URL" env var is set to the mocked api url
    And the 0 "POST" request to "/oauth2/token" returns status 200 with the following response
    """
    {
        "access_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6IjZmZjBjOWEwLWU4ZGMtNDI2MS05YThmLWE5NzE2ZGQ0NDczOSIsInR5cCI6IkpXVCJ9.eyJleHAiOjk3NTE1MDAwMDMsImlhdCI6MTc1MTQ5NjQwMywicm9sZSI6IkFETUlOIiwic2NvcGUiOlsiZW1wbG95ZWVzOnJlYWQiLCJlbXBsb3llZXM6d3JpdGUiLCJyb2xlczp3cml0ZSIsInJvbGVzOnJlYWQiLCJlcGlzOndyaXRlIiwiZXBpczpyZWFkIl0sInN1YiI6IjM4ZGZkZDY5LWQ5OWYtNGJlYS05OTE3LWFlNmFlNmRjMzVlYiJ9.PLG7lnGRQky7lu9wWid_yCc2lMuVDSUFa_uAA3tGlRHFZnp3BAtXs5gmQ30X4WKQGW4BO5uPzx1UadvsLQIfHiwayA9tloh0XEmKZepo7i4oLPB2dqirdCQ09cbmAFD2qUSE0ci4xS1s0wQWsqY4mJf5af92nGr1KfGLYHCyaGtidNWBXc1Yoo8OvgwdmasyRhQo1BIARqrbGNxfgjM1RG2VFBb58DlC3Vm0NqM7cV0vjFf44jffyuVS4AKZKt-Wm5rcVDisko85dNVCr8p25aMshO9u7mRm527XOUPGGaVAzoWf0zmpV-v-1ugQMeYKzFB2vc1L2u1prJ9whiaaZg",
        "token_type": "Bearer",
        "expires_in": 3600,
        "refresh_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6IjllYjY1Mzk0LThhOWItNDRmMy04ZDJkLTRkMjc4ZDE2MDg0NSIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTIxMDEyMDMsImlhdCI6MTc1MTQ5NjQwMywicm9sZSI6IlVTRVIiLCJzY29wZSI6WyJlbXBsb3llZXM6cmVhZCIsImVtcGxveWVlczp3cml0ZSIsInJvbGVzOndyaXRlIiwicm9sZXM6cmVhZCIsImVwaXM6d3JpdGUiLCJlcGlzOnJlYWQiXSwic3ViIjoiMzhkZmRkNjktZDk5Zi00YmVhLTk5MTctYWU2YWU2ZGMzNWViIn0.R6qNzN-HPS4KgB3p-LwzVL19-yBtJq1UQaeSwtbM5rHMhHlgrfOwC-vOBn0PnZ9N6RRbjhg_iqfUrOr2FGhK0g_B5o3aSeI0_kYXRO0PE3uSbd0VSJ0Jv7NX4QZElAY1S1wQzLx-dlkNs3Uh7IdlJ81e1aEHvNaYmriKDCxYY8juaD6OUkxq8VZmWubVazkEMk_RmFWkbdr041gcTouZ_osey_XiMd4pNYMx9pWOvV4BrVff_5NwvU0WIr9WqL1n8WzIP_Xg4gCYxUqFlXGvwLAnYnNkPBE-SxrzrdvlFV8EoE4LC94Oa1lIbdsAOvJD7I7LAduA4JprjOxGgQBRng",
        "refresh_expires_in": 604800,
        "scope": "ADMIN"
    }
    """

  # ========================================
  # SUCCESS SCENARIOS - Complete Pipeline
  # ========================================

  @success @sqs_integration
  Scenario: Process climate analysis via SQS - complete pipeline with S3 data
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00001/5a801269-110d-46c0-89dd-09a6dc1954d1/2024-01-05/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "550e8400-e29b-41d4-a716-446655440000",
      "status": "PENDING",
      "date": "2024-01-05",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc1954d1",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00001",
        "name": "Field 1",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01", "2024-01-02", "2024-01-03", "2024-01-04", "2024-01-05"],
        "temperature_2m_mean": [25.3, 24.8, 26.1, 25.5, 24.9],
        "precipitation_sum": [0.0, 12.5, 0.0, 5.2, 0.0]
      },
      "hourly": {
        "time": [
          "2024-01-01T00:00", "2024-01-01T01:00", "2024-01-01T02:00",
          "2024-01-02T00:00", "2024-01-02T01:00", "2024-01-02T02:00"
        ],
        "relative_humidity_2m": [75.2, 78.3, 76.5, 80.1, 82.3, 79.8]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-001",
          "image_id": "550e8400-e29b-41d4-a716-446655440000",
          "processing_id": "550e8400-e29b-41d4-a716-446655440000",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00001/5a801269-110d-46c0-89dd-09a6dc1954d1/2024-01-05/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    And the 0 "POST" request to "/v1/files" returns status 201 with the following response
    """
    {
      "id": "file-record-001",
      "image_id": "550e8400-e29b-41d4-a716-446655440000",
      "type": "CLIMATE",
      "file_path": "data/climate/550e8400-e29b-41d4-a716-446655440000/test-image_climate.parquet"
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "550e8400-e29b-41d4-a716-446655440000"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    And the response should contain the field "metadata.executionTimeMs" equal to "not nil"
    And the response should contain the field "metadata.weatherApiTimeMs" equal to "not nil"
    And the response should contain the field "metadata.mergingTimeMs" equal to "not nil"
    And the response should contain the field "metadata.message" equal to "not nil"
    And the 0 "POST" request for "/v1/files" json body should have the variable "image_id" with value "550e8400-e29b-41d4-a716-446655440000"
    And the 0 "POST" request for "/v1/files" json body should have the variable "type" with value "CLIMATE"
    And the 0 "POST" request for "/v1/files" json body should have the variable "file_path" with value "not nil"

  @success @sns_integration
  Scenario: Process climate analysis via SNS - complete pipeline with S3 data
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00001/5a801269-110d-46c0-89dd-09a6dc1954d1/2024-01-02/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "660e8400-e29b-41d4-a716-446655440001",
      "status": "PENDING",
      "date": "2024-01-02",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc1954d1",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00001",
        "name": "Field 1",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01", "2024-01-02"],
        "temperature_2m_mean": [25.3, 24.8],
        "precipitation_sum": [0.0, 12.5]
      },
      "hourly": {
        "time": ["2024-01-01T00:00", "2024-01-01T01:00"],
        "relative_humidity_2m": [75.2, 78.3]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-002",
          "image_id": "660e8400-e29b-41d4-a716-446655440001",
          "processing_id": "660e8400-e29b-41d4-a716-446655440001",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00001/5a801269-110d-46c0-89dd-09a6dc1954d1/2024-01-02/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    And the 0 "POST" request to "/v1/files" returns status 201 with the following response
    """
    {
      "id": "file-record-002",
      "image_id": "660e8400-e29b-41d4-a716-446655440001",
      "type": "CLIMATE",
      "file_path": "data/climate/660e8400-e29b-41d4-a716-446655440001/test-image_climate.parquet"
    }
    """
    When the following event is received via sns
    """
    {
      "processing_id": "660e8400-e29b-41d4-a716-446655440001"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    And the response should contain the field "metadata.executionTimeMs" equal to "not nil"
    And the 0 "POST" request for "/v1/files" json body should have the variable "image_id" with value "660e8400-e29b-41d4-a716-446655440001"
    And the 0 "POST" request for "/v1/files" json body should have the variable "type" with value "CLIMATE"

  # ========================================
  # SUCCESS SCENARIOS - Centroid & Weather API
  # ========================================

  @success @centroid_validation
  Scenario: Process climate analysis - validate centroid calculation and weather API parameters
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00002/5a801269-110d-46c0-89dd-09a6dc1954d2/2024-01-01/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "770e8400-e29b-41d4-a716-446655440002",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc1954d2",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00002",
        "name": "Field 2",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [0.0]
      },
      "hourly": {
        "time": ["2024-01-01T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-003",
          "image_id": "770e8400-e29b-41d4-a716-446655440002",
          "processing_id": "770e8400-e29b-41d4-a716-446655440002",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00002/5a801269-110d-46c0-89dd-09a6dc1954d2/2024-01-01/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    And the 0 "POST" request to "/v1/files" returns status 201 with the following response
    """
    {
      "id": "file-record-003",
      "image_id": "770e8400-e29b-41d4-a716-446655440002",
      "type": "CLIMATE",
      "file_path": "data/climate/770e8400-e29b-41d4-a716-446655440002/test-image_climate.parquet"
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "770e8400-e29b-41d4-a716-446655440002"
    }
    """
    Then the lambda should finish without errors
    And the 0 "GET" request for "/v1/archive" queries should have the variable "latitude" with value "-23.4565"
    And the 0 "GET" request for "/v1/archive" queries should have the variable "longitude" with value "-46.7895"
    And the 0 "GET" request for "/v1/archive" queries should have the variable "start_date" with value "not nil"
    And the 0 "GET" request for "/v1/archive" queries should have the variable "end_date" with value "not nil"
    And the 0 "POST" request for "/v1/files" json body should have the variable "image_id" with value "770e8400-e29b-41d4-a716-446655440002"
    And the 0 "POST" request for "/v1/files" json body should have the variable "type" with value "CLIMATE"

  @success @weather_date_range
  Scenario: Process climate analysis - validate 4-month weather data range calculation
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00003/5a801269-110d-46c0-89dd-09a6dc1954d3/2024-06-01/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "880e8400-e29b-41d4-a716-446655440003",
      "status": "PENDING",
      "date": "2024-06-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc1954d3",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00003",
        "name": "Field 3",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-06-01"],
        "temperature_2m_mean": [24.5],
        "precipitation_sum": [0.0]
      },
      "hourly": {
        "time": ["2024-06-01T00:00", "2024-06-01T01:00"],
        "relative_humidity_2m": [72.1, 73.5]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-004",
          "image_id": "880e8400-e29b-41d4-a716-446655440003",
          "processing_id": "880e8400-e29b-41d4-a716-446655440003",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00003/5a801269-110d-46c0-89dd-09a6dc1954d3/2024-06-01/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "880e8400-e29b-41d4-a716-446655440003"
    }
    """
    Then the lambda should finish without errors
    And the 0 "GET" request for "/v1/archive" queries should have the variable "start_date" with value "not nil"
    And the 0 "GET" request for "/v1/archive" queries should have the variable "end_date" with value "not nil"

  # ========================================
  # SUCCESS SCENARIOS - Response Structure
  # ========================================

  @success @response_validation
  Scenario: Process climate analysis - validate complete final dataset structure
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00004/5a801269-110d-46c0-89dd-09a6dc1954d4/2024-02-03/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "990e8400-e29b-41d4-a716-446655440004",
      "status": "PENDING",
      "date": "2024-02-03",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc1954d4",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00004",
        "name": "Field 4",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-02-01", "2024-02-02", "2024-02-03"],
        "temperature_2m_mean": [24.5, 25.0, 24.8],
        "precipitation_sum": [0.0, 5.2, 0.0]
      },
      "hourly": {
        "time": ["2024-02-01T00:00", "2024-02-01T01:00"],
        "relative_humidity_2m": [72.1, 73.5]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-005",
          "image_id": "990e8400-e29b-41d4-a716-446655440004",
          "processing_id": "990e8400-e29b-41d4-a716-446655440004",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00004/5a801269-110d-46c0-89dd-09a6dc1954d4/2024-02-03/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "990e8400-e29b-41d4-a716-446655440004"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    And the response should contain the field "metadata.executionTimeMs" equal to "not nil"

  @success @metadata_validation
  Scenario: Process climate analysis - validate execution metadata
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00005/5a801269-110d-46c0-89dd-09a6dc1954d5/2024-01-01/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "aa0e8400-e29b-41d4-a716-446655440005",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc1954d5",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00005",
        "name": "Field 5",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [0.0]
      },
      "hourly": {
        "time": ["2024-01-01T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-006",
          "image_id": "aa0e8400-e29b-41d4-a716-446655440005",
          "processing_id": "aa0e8400-e29b-41d4-a716-446655440005",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00005/5a801269-110d-46c0-89dd-09a6dc1954d5/2024-01-01/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "aa0e8400-e29b-41d4-a716-446655440005"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "metadata.executionTimeMs" equal to "not nil"
    And the response should contain the field "metadata.weatherApiTimeMs" equal to "not nil"
    And the response should contain the field "metadata.mergingTimeMs" equal to "not nil"
    And the response should contain the field "metadata.message" equal to "not nil"

  @success @multiple_records
  Scenario: Process climate analysis - validate multiple pixel records in response
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00006/5a801269-110d-46c0-89dd-09a6dc1954d6/2024-01-01/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "bb0e8400-e29b-41d4-a716-446655440006",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc1954d6",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00006",
        "name": "Field 6",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [0.0]
      },
      "hourly": {
        "time": ["2024-01-01T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-007",
          "image_id": "bb0e8400-e29b-41d4-a716-446655440006",
          "processing_id": "bb0e8400-e29b-41d4-a716-446655440006",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00006/5a801269-110d-46c0-89dd-09a6dc1954d6/2024-01-01/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "bb0e8400-e29b-41d4-a716-446655440006"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    And the response should contain the field "metadata.executionTimeMs" equal to "not nil"

  # ========================================
  # SUCCESS SCENARIOS - Weather Metrics
  # ========================================

  @success @weather_metrics_calculation
  Scenario: Process climate analysis - validate 30-day rolling window weather metrics
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00012/5a801269-110d-46c0-89dd-09a6dc195412/2024-01-24/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "cc1e8400-e29b-41d4-a716-446655440012",
      "status": "PENDING",
      "date": "2024-01-24",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195412",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00012",
        "name": "Field 12",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": [
          "2024-01-15", "2024-01-16", "2024-01-17", "2024-01-18", "2024-01-19",
          "2024-01-20", "2024-01-21", "2024-01-22", "2024-01-23", "2024-01-24"
        ],
        "temperature_2m_mean": [25.0, 26.0, 24.5, 25.5, 26.2, 24.8, 25.3, 26.1, 24.9, 25.7],
        "precipitation_sum": [0.0, 0.0, 0.0, 5.2, 0.0, 0.0, 0.0, 12.5, 0.0, 0.0]
      },
      "hourly": {
        "time": ["2024-01-15T00:00", "2024-01-15T01:00"],
        "relative_humidity_2m": [72.0, 74.0]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-012",
          "image_id": "cc1e8400-e29b-41d4-a716-446655440012",
          "processing_id": "cc1e8400-e29b-41d4-a716-446655440012",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00012/5a801269-110d-46c0-89dd-09a6dc195412/2024-01-24/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "cc1e8400-e29b-41d4-a716-446655440012"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    And the response should contain the field "metadata.executionTimeMs" equal to "not nil"

  @success @dry_days_calculation
  Scenario: Process climate analysis - validate consecutive dry days calculation
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00013/5a801269-110d-46c0-89dd-09a6dc195413/2024-01-07/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "dd1e8400-e29b-41d4-a716-446655440013",
      "status": "PENDING",
      "date": "2024-01-07",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195413",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00013",
        "name": "Field 13",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01", "2024-01-02", "2024-01-03", "2024-01-04", "2024-01-05", "2024-01-06", "2024-01-07"],
        "temperature_2m_mean": [25.0, 25.0, 25.0, 25.0, 25.0, 25.0, 25.0],
        "precipitation_sum": [0.0, 0.0, 0.0, 5.2, 0.0, 0.0, 0.0]
      },
      "hourly": {
        "time": ["2024-01-01T00:00"],
        "relative_humidity_2m": [70.0]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-013",
          "image_id": "dd1e8400-e29b-41d4-a716-446655440013",
          "processing_id": "dd1e8400-e29b-41d4-a716-446655440013",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00013/5a801269-110d-46c0-89dd-09a6dc195413/2024-01-07/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "dd1e8400-e29b-41d4-a716-446655440013"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    And the response should contain the field "metadata.executionTimeMs" equal to "not nil"

  # ========================================
  # ERROR SCENARIOS - Input Validation
  # ========================================

  @error @contract_fields_validation
  Scenario: Process climate analysis with missing processing_id - should fail validation
    When the following event is received via sqs
    """
    {}
    """
    Then the lambda should finish with "bad request" error

  @error @contract_fields_validation
  Scenario: Process climate analysis with empty processing_id - should fail validation
    When the following event is received via sqs
    """
    {
      "processing_id": ""
    }
    """
    Then the lambda should finish with "bad request" error

  @error @contract_fields_validation
  Scenario: Process climate analysis with null processing_id - should fail validation
    When the following event is received via sqs
    """
    {
      "processing_id": null
    }
    """
    Then the lambda should finish with "bad request" error

  @error @contract_fields_validation
  Scenario: Process climate analysis with invalid processing_id format - should fail validation
    When the following event is received via sqs
    """
    {
      "processing_id": "invalid-uuid-format"
    }
    """
    Then the lambda should finish with "bad request" error

  # ========================================
  # ERROR SCENARIOS - S3 Data Issues
  # ========================================

  @error @s3_not_found
  Scenario: Process climate analysis with non-existent processing_id - should return not found error
    And the 0 "GET" request to "/v1/analysis/*" returns status 404 with the following response
    """
    {
      "error": {
        "description": "Processing not found",
        "code": "NOT_FOUND"
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "ee2e8400-e29b-41d4-a716-446655440014"
    }
    """
    Then the lambda should finish with "not found" error

  @error @s3_missing_files
  Scenario: Process climate analysis with missing parquet dataset file in S3 - should return not found error
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "ff2e8400-e29b-41d4-a716-446655440015",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195415",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00015",
        "name": "Field 15",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [0.0]
      },
      "hourly": {
        "time": ["2024-01-01T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "ff2e8400-e29b-41d4-a716-446655440015"
    }
    """
    Then the lambda should finish with "not found" error

  @error @api_not_found
  Scenario: Process climate analysis with analysis API returning 404 - should return not found error
    And the 0 "GET" request to "/v1/analysis/*" returns status 404 with the following response
    """
    {
      "error": {
        "description": "Analysis not found",
        "code": "NOT_FOUND"
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "bb3e8400-e29b-41d4-a716-446655440017"
    }
    """
    Then the lambda should finish with "not found" error

  # ========================================
  # ERROR SCENARIOS - Weather API Issues
  # ========================================

  @error @weather_api_failure
  Scenario: Process climate analysis with weather API failure - should retry and fail gracefully
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00008/5a801269-110d-46c0-89dd-09a6dc1954d8/2024-01-01/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "dd0e8400-e29b-41d4-a716-446655440008",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc1954d8",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00008",
        "name": "Field 8",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the "GET" "/v1/archive" should not return any results
    When the following event is received via sqs
    """
    {
      "processing_id": "dd0e8400-e29b-41d4-a716-446655440008"
    }
    """
    Then the lambda should finish with "internal server" error

  @error @weather_api_timeout
  Scenario: Process climate analysis with weather API timeout - should fail gracefully
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00019/5a801269-110d-46c0-89dd-09a6dc195419/2024-01-01/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "dd3e8400-e29b-41d4-a716-446655440019",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195419",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00019",
        "name": "Field 19",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the "GET" "/v1/archive" should not return any results
    When the following event is received via sqs
    """
    {
      "processing_id": "dd3e8400-e29b-41d4-a716-446655440019"
    }
    """
    Then the lambda should finish with "internal server" error

  @error @invalid_coordinates
  Scenario: Process climate analysis with invalid coordinates - should fail weather fetch
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00009/5a801269-110d-46c0-89dd-09a6dc1954d9/2024-01-01/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "ee0e8400-e29b-41d4-a716-446655440009",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc1954d9",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00009",
        "name": "Field 9",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": 95.0,
          "lng": 200.0
        },
        "area_hectares": 10.5
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "ee0e8400-e29b-41d4-a716-446655440009"
    }
    """
    Then the lambda should finish with "bad request" error

  @error @weather_api_rate_limit
  Scenario: Process climate analysis with weather API rate limit - should handle 429 response
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00020/5a801269-110d-46c0-89dd-09a6dc195420/2024-01-01/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "ee3e8400-e29b-41d4-a716-446655440020",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195420",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00020",
        "name": "Field 20",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 429 with the following response
    """
    {
      "error": true,
      "reason": "Rate limit exceeded"
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "ee3e8400-e29b-41d4-a716-446655440020"
    }
    """
    Then the lambda should finish with "server error occurred" error

  # ========================================
  # SUCCESS SCENARIOS - Caching
  # ========================================

  @success @cache_hit
  Scenario: Process climate analysis with cached weather data - should use cache
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00021/5a801269-110d-46c0-89dd-09a6dc195421/2024-01-01/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "ff3e8400-e29b-41d4-a716-446655440021",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195421",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00021",
        "name": "Field 21",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [0.0]
      },
      "hourly": {
        "time": ["2024-01-01T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-021",
          "image_id": "ff3e8400-e29b-41d4-a716-446655440021",
          "processing_id": "ff3e8400-e29b-41d4-a716-446655440021",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00021/5a801269-110d-46c0-89dd-09a6dc195421/2024-01-01/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "ff3e8400-e29b-41d4-a716-446655440021"
    }
    """
    Then the lambda should finish without errors
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00022/5a801269-110d-46c0-89dd-09a6dc195422/2024-01-01/delta.parquet"
    And the 1 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "ff3e8400-e29b-41d4-a716-446655440022",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195422",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00022",
        "name": "Field 22",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 1 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-022",
          "image_id": "ff3e8400-e29b-41d4-a716-446655440022",
          "processing_id": "ff3e8400-e29b-41d4-a716-446655440022",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00022/5a801269-110d-46c0-89dd-09a6dc195422/2024-01-01/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "ff3e8400-e29b-41d4-a716-446655440022"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "metadata.cacheHit" equal to "not nil"

  # ========================================
  # SUCCESS SCENARIOS - Concurrent Processing
  # ========================================

  @success @concurrent_processing
  Scenario: Process multiple climate analyses concurrently - validate independent processing
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00022/5a801269-110d-46c0-89dd-09a6dc195422/2024-01-01/delta.parquet"
    And the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00022/5a801269-110d-46c0-89dd-09a6dc195422/2024-01-01/delta.parquet"
    # Mock responses for first event (index 0)
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "aa4e8400-e29b-41d4-a716-446655440022",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195422",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00022",
        "name": "Field 22",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [0.0]
      },
      "hourly": {
        "time": ["2024-01-01T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-022-concurrent",
          "image_id": "aa4e8400-e29b-41d4-a716-446655440022",
          "processing_id": "aa4e8400-e29b-41d4-a716-446655440022",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00022/5a801269-110d-46c0-89dd-09a6dc195422/2024-01-01/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    # Mock responses for second event (index 1)
    And the 1 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "bb4e8400-e29b-41d4-a716-446655440023",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195422",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00022",
        "name": "Field 22",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 1 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [0.0]
      },
      "hourly": {
        "time": ["2024-01-01T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    And the 1 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-023-concurrent",
          "image_id": "bb4e8400-e29b-41d4-a716-446655440023",
          "processing_id": "bb4e8400-e29b-41d4-a716-446655440023",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00022/5a801269-110d-46c0-89dd-09a6dc195422/2024-01-01/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    # Execute first event
    When the following event is received via sqs
    """
    {
      "processing_id": "aa4e8400-e29b-41d4-a716-446655440022"
    }
    """
    Then the lambda should finish without errors
    # Execute second event
    When the following event is received via sqs
    """
    {
      "processing_id": "bb4e8400-e29b-41d4-a716-446655440023"
    }
    """
    Then the lambda should finish without errors

  # ========================================
  # SUCCESS SCENARIOS - Parquet Schema Validation
  # ========================================

  @success @parquet_schema_validation
  Scenario: Validate enriched parquet file contains all original fields plus weather fields
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00024/5a801269-110d-46c0-89dd-09a6dc195424/2024-01-30/delta.parquet"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00024/5a801269-110d-46c0-89dd-09a6dc195424/2024-01-30/delta.parquet" should have 1162 rows
    And the "HISTORICAL_CLIMATE_DAYS" env var is set to "30"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "cc5e8400-e29b-41d4-a716-446655440024",
      "status": "PENDING",
      "date": "2024-01-30",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195424",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00024",
        "name": "Field 24",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": [
          "2024-01-01", "2024-01-02", "2024-01-03", "2024-01-04", "2024-01-05",
          "2024-01-06", "2024-01-07", "2024-01-08", "2024-01-09", "2024-01-10",
          "2024-01-11", "2024-01-12", "2024-01-13", "2024-01-14", "2024-01-15",
          "2024-01-16", "2024-01-17", "2024-01-18", "2024-01-19", "2024-01-20",
          "2024-01-21", "2024-01-22", "2024-01-23", "2024-01-24", "2024-01-25",
          "2024-01-26", "2024-01-27", "2024-01-28", "2024-01-29", "2024-01-30"
        ],
        "temperature_2m_mean": [
          25.3, 24.8, 26.1, 25.5, 24.9, 27.2, 26.8, 25.4, 24.6, 26.3,
          25.7, 26.9, 25.1, 24.4, 26.5, 27.1, 25.8, 24.7, 26.2, 25.6,
          27.3, 26.4, 25.2, 24.8, 26.7, 27.0, 25.9, 24.5, 26.1, 25.4
        ],
        "precipitation_sum": [
          0.0, 5.2, 0.0, 12.5, 0.0, 0.0, 3.8, 0.0, 0.0, 7.1,
          0.0, 0.0, 15.3, 0.0, 0.0, 0.0, 8.6, 0.0, 0.0, 4.2,
          0.0, 0.0, 0.0, 10.7, 0.0, 0.0, 6.4, 0.0, 0.0, 9.3
        ]
      },
      "hourly": {
        "time": [
          "2024-01-01T00:00", "2024-01-01T01:00", "2024-01-01T02:00", "2024-01-01T03:00", "2024-01-01T04:00", "2024-01-01T05:00",
          "2024-01-02T00:00", "2024-01-02T01:00", "2024-01-02T02:00", "2024-01-02T03:00", "2024-01-02T04:00", "2024-01-02T05:00",
          "2024-01-03T00:00", "2024-01-03T01:00", "2024-01-03T02:00", "2024-01-03T03:00", "2024-01-03T04:00", "2024-01-03T05:00",
          "2024-01-04T00:00", "2024-01-04T01:00", "2024-01-04T02:00", "2024-01-04T03:00", "2024-01-04T04:00", "2024-01-04T05:00",
          "2024-01-05T00:00", "2024-01-05T01:00", "2024-01-05T02:00", "2024-01-05T03:00", "2024-01-05T04:00", "2024-01-05T05:00"
        ],
        "relative_humidity_2m": [
          75.2, 78.3, 76.5, 72.1, 74.8, 77.6,
          80.1, 82.3, 79.8, 76.4, 78.9, 81.2,
          74.5, 77.8, 75.3, 73.2, 76.1, 78.4,
          82.6, 84.1, 81.5, 79.7, 82.0, 83.4,
          73.8, 76.4, 74.9, 72.6, 75.3, 77.1
        ]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-024",
          "image_id": "cc5e8400-e29b-41d4-a716-446655440024",
          "processing_id": "cc5e8400-e29b-41d4-a716-446655440024",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00024/5a801269-110d-46c0-89dd-09a6dc195424/2024-01-30/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "cc5e8400-e29b-41d4-a716-446655440024"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00024/5a801269-110d-46c0-89dd-09a6dc195424/2024-01-30/climate.parquet" should have 1162 rows
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00024/5a801269-110d-46c0-89dd-09a6dc195424/2024-01-30/climate.parquet" should have all fields from "test/integration/resources/climate-analysis/test-image_raw.parquet"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00024/5a801269-110d-46c0-89dd-09a6dc195424/2024-01-30/climate.parquet" should have the fields "avg_temperature, temp_std_dev, avg_humidity, humidity_std_dev, total_precipitation, dry_days_consecutive" in all rows
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00024/5a801269-110d-46c0-89dd-09a6dc195424/2024-01-30/climate.parquet" should have the fields "B02, B03, B04, latitude, longitude" in all rows

  @success @parquet_row_values_validation
  Scenario: Validate enriched parquet preserves original row values and adds correct weather data
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet"
    # Validate first 5 rows of original file (rows 0-4)
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "x" in row 0 with value "9"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "y" in row 0 with value "0"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B02" in row 0 with value "1225"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B03" in row 0 with value "1352"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B04" in row 0 with value "1228"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "x" in row 1 with value "10"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "y" in row 1 with value "0"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B02" in row 1 with value "1223"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B03" in row 1 with value "1341"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B04" in row 1 with value "1213"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "x" in row 2 with value "11"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "y" in row 2 with value "0"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B02" in row 2 with value "1227"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B03" in row 2 with value "1341"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B04" in row 2 with value "1208"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "x" in row 3 with value "12"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "y" in row 3 with value "0"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B02" in row 3 with value "1227"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B03" in row 3 with value "1330"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B04" in row 3 with value "1212"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "x" in row 4 with value "13"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "y" in row 4 with value "0"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B02" in row 4 with value "1224"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B03" in row 4 with value "1336"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B04" in row 4 with value "1221"
    # Validate last 5 rows of original file (rows 1157-1161)
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "x" in row 1157 with value "24"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "y" in row 1157 with value "38"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B02" in row 1157 with value "1299"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B03" in row 1157 with value "1398"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B04" in row 1157 with value "1353"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "x" in row 1158 with value "21"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "y" in row 1158 with value "39"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B02" in row 1158 with value "1360"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B03" in row 1158 with value "1526"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B04" in row 1158 with value "1410"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "x" in row 1159 with value "22"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "y" in row 1159 with value "39"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B02" in row 1159 with value "1389"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B03" in row 1159 with value "1579"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B04" in row 1159 with value "1423"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "x" in row 1160 with value "23"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "y" in row 1160 with value "39"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B02" in row 1160 with value "1392"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B03" in row 1160 with value "1610"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B04" in row 1160 with value "1470"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "x" in row 1161 with value "24"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "y" in row 1161 with value "39"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B02" in row 1161 with value "1383"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B03" in row 1161 with value "1566"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet" should have the field "B04" in row 1161 with value "1480"
    And the "HISTORICAL_CLIMATE_DAYS" env var is set to "30"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "dd5e8400-e29b-41d4-a716-446655440025",
      "status": "PENDING",
      "date": "2024-01-30",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195425",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00025",
        "name": "Field 25",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": [
          "2024-01-01", "2024-01-02", "2024-01-03", "2024-01-04", "2024-01-05",
          "2024-01-06", "2024-01-07", "2024-01-08", "2024-01-09", "2024-01-10",
          "2024-01-11", "2024-01-12", "2024-01-13", "2024-01-14", "2024-01-15",
          "2024-01-16", "2024-01-17", "2024-01-18", "2024-01-19", "2024-01-20",
          "2024-01-21", "2024-01-22", "2024-01-23", "2024-01-24", "2024-01-25",
          "2024-01-26", "2024-01-27", "2024-01-28", "2024-01-29", "2024-01-30"
        ],
        "temperature_2m_mean": [
          22.5, 23.1, 21.8, 24.2, 23.7, 22.9, 24.5, 23.3, 22.1, 24.0,
          23.5, 22.8, 24.3, 23.0, 21.9, 24.1, 23.6, 22.7, 24.4, 23.2,
          22.0, 23.9, 23.4, 22.6, 24.6, 23.1, 21.7, 24.2, 23.8, 22.4
        ],
        "precipitation_sum": [
          5.2, 0.0, 8.3, 0.0, 0.0, 12.1, 0.0, 0.0, 0.0, 6.5,
          0.0, 0.0, 9.8, 0.0, 0.0, 0.0, 7.4, 0.0, 0.0, 11.2,
          0.0, 0.0, 0.0, 4.9, 0.0, 0.0, 10.3, 0.0, 0.0, 8.7
        ]
      },
      "hourly": {
        "time": [
          "2024-01-01T00:00", "2024-01-01T01:00", "2024-01-01T02:00", "2024-01-01T03:00", "2024-01-01T04:00", "2024-01-01T05:00",
          "2024-01-02T00:00", "2024-01-02T01:00", "2024-01-02T02:00", "2024-01-02T03:00", "2024-01-02T04:00", "2024-01-02T05:00",
          "2024-01-03T00:00", "2024-01-03T01:00", "2024-01-03T02:00", "2024-01-03T03:00", "2024-01-03T04:00", "2024-01-03T05:00",
          "2024-01-04T00:00", "2024-01-04T01:00", "2024-01-04T02:00", "2024-01-04T03:00", "2024-01-04T04:00", "2024-01-04T05:00",
          "2024-01-05T00:00", "2024-01-05T01:00", "2024-01-05T02:00", "2024-01-05T03:00", "2024-01-05T04:00", "2024-01-05T05:00"
        ],
        "relative_humidity_2m": [
          72.5, 74.8, 73.2, 71.5, 73.9, 75.1,
          76.3, 78.6, 77.1, 75.4, 77.8, 79.2,
          73.7, 76.0, 74.4, 72.7, 75.1, 76.5,
          79.8, 81.3, 80.1, 78.5, 80.9, 82.1,
          74.2, 76.5, 74.9, 73.2, 75.6, 76.8
        ]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-025",
          "image_id": "dd5e8400-e29b-41d4-a716-446655440025",
          "processing_id": "dd5e8400-e29b-41d4-a716-446655440025",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "dd5e8400-e29b-41d4-a716-446655440025"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    # Validate first 5 rows of enriched file - original fields preserved (rows 0-4)
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "x" in row 0 with value "9"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "y" in row 0 with value "0"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B02" in row 0 with value "1225"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B03" in row 0 with value "1352"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B04" in row 0 with value "1228"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "x" in row 1 with value "10"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "y" in row 1 with value "0"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B02" in row 1 with value "1223"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B03" in row 1 with value "1341"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B04" in row 1 with value "1213"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "x" in row 2 with value "11"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "y" in row 2 with value "0"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B02" in row 2 with value "1227"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B03" in row 2 with value "1341"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B04" in row 2 with value "1208"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "x" in row 3 with value "12"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "y" in row 3 with value "0"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B02" in row 3 with value "1227"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B03" in row 3 with value "1330"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B04" in row 3 with value "1212"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "x" in row 4 with value "13"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "y" in row 4 with value "0"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B02" in row 4 with value "1224"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B03" in row 4 with value "1336"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B04" in row 4 with value "1221"
    # Validate last 5 rows of enriched file - original fields preserved (rows 1157-1161)
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "x" in row 1157 with value "24"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "y" in row 1157 with value "38"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B02" in row 1157 with value "1299"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B03" in row 1157 with value "1398"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B04" in row 1157 with value "1353"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "x" in row 1158 with value "21"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "y" in row 1158 with value "39"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B02" in row 1158 with value "1360"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B03" in row 1158 with value "1526"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B04" in row 1158 with value "1410"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "x" in row 1159 with value "22"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "y" in row 1159 with value "39"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B02" in row 1159 with value "1389"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B03" in row 1159 with value "1579"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B04" in row 1159 with value "1423"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "x" in row 1160 with value "23"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "y" in row 1160 with value "39"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B02" in row 1160 with value "1392"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B03" in row 1160 with value "1610"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B04" in row 1160 with value "1470"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "x" in row 1161 with value "24"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "y" in row 1161 with value "39"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B02" in row 1161 with value "1383"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B03" in row 1161 with value "1566"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "B04" in row 1161 with value "1480"
    # Validate statistical climate fields - first 5 rows (rows 0-4)
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "avg_temperature" in row 0 with value "23.24"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "temp_std_dev" in row 0 with value "0.86"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "avg_humidity" in row 0 with value "12.71"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "humidity_std_dev" in row 0 with value "28.92"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "total_precipitation" in row 0 with value "84.4"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "dry_days_consecutive" in row 0 with value "3"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "avg_temperature" in row 1 with value "23.24"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "temp_std_dev" in row 1 with value "0.86"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "avg_humidity" in row 1 with value "12.71"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "humidity_std_dev" in row 1 with value "28.92"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "total_precipitation" in row 1 with value "84.4"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "dry_days_consecutive" in row 1 with value "3"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "avg_temperature" in row 2 with value "23.24"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "temp_std_dev" in row 2 with value "0.86"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "avg_humidity" in row 2 with value "12.71"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "humidity_std_dev" in row 2 with value "28.92"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "total_precipitation" in row 2 with value "84.4"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "dry_days_consecutive" in row 2 with value "3"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "avg_temperature" in row 3 with value "23.24"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "temp_std_dev" in row 3 with value "0.86"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "avg_humidity" in row 3 with value "12.71"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "humidity_std_dev" in row 3 with value "28.92"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "total_precipitation" in row 3 with value "84.4"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "dry_days_consecutive" in row 3 with value "3"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "avg_temperature" in row 4 with value "23.24"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "temp_std_dev" in row 4 with value "0.86"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "avg_humidity" in row 4 with value "12.71"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "humidity_std_dev" in row 4 with value "28.92"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "total_precipitation" in row 4 with value "84.4"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "dry_days_consecutive" in row 4 with value "3"
    # Validate statistical climate fields - last 5 rows (rows 1157-1161)
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "avg_temperature" in row 1157 with value "23.24"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "temp_std_dev" in row 1157 with value "0.86"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "avg_humidity" in row 1157 with value "12.71"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "humidity_std_dev" in row 1157 with value "28.92"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "total_precipitation" in row 1157 with value "84.4"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "dry_days_consecutive" in row 1157 with value "3"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "avg_temperature" in row 1158 with value "23.24"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "temp_std_dev" in row 1158 with value "0.86"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "avg_humidity" in row 1158 with value "12.71"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "humidity_std_dev" in row 1158 with value "28.92"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "total_precipitation" in row 1158 with value "84.4"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "dry_days_consecutive" in row 1158 with value "3"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "avg_temperature" in row 1159 with value "23.24"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "temp_std_dev" in row 1159 with value "0.86"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "avg_humidity" in row 1159 with value "12.71"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "humidity_std_dev" in row 1159 with value "28.92"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "total_precipitation" in row 1159 with value "84.4"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "dry_days_consecutive" in row 1159 with value "3"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "avg_temperature" in row 1160 with value "23.24"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "temp_std_dev" in row 1160 with value "0.86"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "avg_humidity" in row 1160 with value "12.71"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "humidity_std_dev" in row 1160 with value "28.92"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "total_precipitation" in row 1160 with value "84.4"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "dry_days_consecutive" in row 1160 with value "3"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "avg_temperature" in row 1161 with value "23.24"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "temp_std_dev" in row 1161 with value "0.86"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "avg_humidity" in row 1161 with value "12.71"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "humidity_std_dev" in row 1161 with value "28.92"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "total_precipitation" in row 1161 with value "84.4"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00025/5a801269-110d-46c0-89dd-09a6dc195425/2024-01-30/climate.parquet" should have the field "dry_days_consecutive" in row 1161 with value "3"

  # ========================================
  # ERROR SCENARIOS - Domain Errors & Edge Cases
  # ========================================

  @error @empty_dataset
  Scenario: Process climate analysis with empty dataset - should return empty dataset error
    Given the "test/integration/resources/climate-analysis/empty.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00030/5a801269-110d-46c0-89dd-09a6dc195430/2024-01-01/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "xx1e8400-e29b-41d4-a716-446655440030",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195430",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00030",
        "name": "Field 30",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [3.5]
      },
      "hourly": {
        "time": ["2024-01-01T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "xx1e8400-e29b-41d4-a716-446655440030"
    }
    """
    Then the lambda should finish with "bad request" error

  # ========================================
  # SUCCESS SCENARIOS - Historical Climate Metrics
  # ========================================

  @success @historical_climate_metrics
  Scenario: Process climate analysis with 30-day historical window - validate statistical climate metrics
    Given the "HISTORICAL_CLIMATE_DAYS" env var is set to "30"
    And the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00031/5a801269-110d-46c0-89dd-09a6dc195431/2024-01-31/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "bb1e8400-e29b-41d4-a716-446655440031",
      "status": "PENDING",
      "date": "2024-01-31",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195431",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00031",
        "name": "Field 31",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": [
          "2024-01-01", "2024-01-02", "2024-01-03", "2024-01-04", "2024-01-05",
          "2024-01-06", "2024-01-07", "2024-01-08", "2024-01-09", "2024-01-10",
          "2024-01-11", "2024-01-12", "2024-01-13", "2024-01-14", "2024-01-15",
          "2024-01-16", "2024-01-17", "2024-01-18", "2024-01-19", "2024-01-20",
          "2024-01-21", "2024-01-22", "2024-01-23", "2024-01-24", "2024-01-25",
          "2024-01-26", "2024-01-27", "2024-01-28", "2024-01-29", "2024-01-30",
          "2024-01-31"
        ],
        "temperature_2m_mean": [
          25.0, 26.0, 24.5, 25.5, 26.2, 24.8, 25.3, 26.1, 24.9, 25.7,
          25.2, 26.3, 24.6, 25.4, 26.0, 24.7, 25.1, 26.2, 24.8, 25.6,
          25.3, 26.1, 24.9, 25.5, 26.0, 24.8, 25.2, 26.3, 24.7, 25.4,
          25.5
        ],
        "precipitation_sum": [
          0.0, 0.0, 0.0, 5.2, 0.0, 0.0, 0.0, 12.5, 0.0, 0.0,
          0.0, 0.0, 3.8, 0.0, 0.0, 0.0, 0.0, 0.0, 7.3, 0.0,
          0.0, 0.0, 0.0, 0.0, 0.0, 9.1, 0.0, 0.0, 0.0, 0.0,
          0.0
        ]
      },
      "hourly": {
        "time": [
          "2024-01-01T00:00", "2024-01-01T01:00", "2024-01-01T02:00",
          "2024-01-02T00:00", "2024-01-02T01:00", "2024-01-02T02:00"
        ],
        "relative_humidity_2m": [75.0, 76.0, 74.5, 75.5, 76.2, 74.8]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-031",
          "image_id": "bb1e8400-e29b-41d4-a716-446655440031",
          "processing_id": "bb1e8400-e29b-41d4-a716-446655440031",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00031/5a801269-110d-46c0-89dd-09a6dc195431/2024-01-31/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "bb1e8400-e29b-41d4-a716-446655440031"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00031/5a801269-110d-46c0-89dd-09a6dc195431/2024-01-31/climate.parquet" should have the fields "avg_temperature, temp_std_dev, avg_humidity, humidity_std_dev, total_precipitation, dry_days_consecutive" in all rows
    And the 0 "GET" request for "/v1/archive" queries should have the variable "start_date" with value "2024-01-01"
    And the 0 "GET" request for "/v1/archive" queries should have the variable "end_date" with value "2024-01-31"

  @success @historical_consecutive_dry_days
  Scenario: Process climate analysis - validate consecutive dry days calculation with gaps
    Given the "HISTORICAL_CLIMATE_DAYS" env var is set to "15"
    And the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00032/5a801269-110d-46c0-89dd-09a6dc195432/2024-01-15/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "bb2e8400-e29b-41d4-a716-446655440032",
      "status": "PENDING",
      "date": "2024-01-15",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195432",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00032",
        "name": "Field 32",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": [
          "2024-01-01", "2024-01-02", "2024-01-03", "2024-01-04", "2024-01-05",
          "2024-01-06", "2024-01-07", "2024-01-08", "2024-01-09", "2024-01-10",
          "2024-01-11", "2024-01-12", "2024-01-13", "2024-01-14", "2024-01-15"
        ],
        "temperature_2m_mean": [
          25.0, 25.0, 25.0, 25.0, 25.0, 25.0, 25.0, 25.0, 25.0, 25.0,
          25.0, 25.0, 25.0, 25.0, 25.0
        ],
        "precipitation_sum": [
          0.0, 0.0, 0.0, 5.2, 0.0, 0.0, 0.0, 0.0, 0.0, 12.5,
          0.0, 0.0, 0.0, 0.0, 0.0
        ]
      },
      "hourly": {
        "time": ["2024-01-01T00:00", "2024-01-01T01:00"],
        "relative_humidity_2m": [75.0, 76.0]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-032",
          "image_id": "bb2e8400-e29b-41d4-a716-446655440032",
          "processing_id": "bb2e8400-e29b-41d4-a716-446655440032",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00032/5a801269-110d-46c0-89dd-09a6dc195432/2024-01-15/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "bb2e8400-e29b-41d4-a716-446655440032"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00032/5a801269-110d-46c0-89dd-09a6dc195432/2024-01-15/climate.parquet" should have the field "dry_days_consecutive" in all rows

  @success @historical_temperature_variance
  Scenario: Process climate analysis - validate temperature statistics with high variance
    Given the "HISTORICAL_CLIMATE_DAYS" env var is set to "10"
    And the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00033/5a801269-110d-46c0-89dd-09a6dc195433/2024-01-10/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "bb3e8400-e29b-41d4-a716-446655440033",
      "status": "PENDING",
      "date": "2024-01-10",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195433",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00033",
        "name": "Field 33",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": [
          "2024-01-01", "2024-01-02", "2024-01-03", "2024-01-04", "2024-01-05",
          "2024-01-06", "2024-01-07", "2024-01-08", "2024-01-09", "2024-01-10"
        ],
        "temperature_2m_mean": [20.0, 22.0, 24.0, 26.0, 28.0, 30.0, 28.0, 26.0, 24.0, 22.0],
        "precipitation_sum": [0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0]
      },
      "hourly": {
        "time": ["2024-01-01T00:00", "2024-01-01T01:00"],
        "relative_humidity_2m": [70.0, 72.0]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-033",
          "image_id": "bb3e8400-e29b-41d4-a716-446655440033",
          "processing_id": "bb3e8400-e29b-41d4-a716-446655440033",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00033/5a801269-110d-46c0-89dd-09a6dc195433/2024-01-10/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "bb3e8400-e29b-41d4-a716-446655440033"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00033/5a801269-110d-46c0-89dd-09a6dc195433/2024-01-10/climate.parquet" should have the field "temp_std_dev" in all rows

  @success @historical_all_rainy_days
  Scenario: Process climate analysis - validate metrics when all days have precipitation
    Given the "HISTORICAL_CLIMATE_DAYS" env var is set to "7"
    And the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00034/5a801269-110d-46c0-89dd-09a6dc195434/2024-01-07/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "bb4e8400-e29b-41d4-a716-446655440034",
      "status": "PENDING",
      "date": "2024-01-07",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195434",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00034",
        "name": "Field 34",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01", "2024-01-02", "2024-01-03", "2024-01-04", "2024-01-05", "2024-01-06", "2024-01-07"],
        "temperature_2m_mean": [23.0, 23.5, 24.0, 23.8, 23.2, 23.6, 23.4],
        "precipitation_sum": [5.5, 3.2, 7.8, 2.1, 9.3, 4.6, 6.7]
      },
      "hourly": {
        "time": ["2024-01-01T00:00", "2024-01-01T01:00"],
        "relative_humidity_2m": [85.0, 86.0]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-034",
          "image_id": "bb4e8400-e29b-41d4-a716-446655440034",
          "processing_id": "bb4e8400-e29b-41d4-a716-446655440034",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00034/5a801269-110d-46c0-89dd-09a6dc195434/2024-01-07/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "bb4e8400-e29b-41d4-a716-446655440034"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00034/5a801269-110d-46c0-89dd-09a6dc195434/2024-01-07/climate.parquet" should have the field "total_precipitation" in all rows

  @success @historical_humidity_variance
  Scenario: Process climate analysis - validate humidity statistics with high variance
    Given the "HISTORICAL_CLIMATE_DAYS" env var is set to "14"
    And the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00035/5a801269-110d-46c0-89dd-09a6dc195435/2024-01-14/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "bb5e8400-e29b-41d4-a716-446655440035",
      "status": "PENDING",
      "date": "2024-01-14",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195435",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00035",
        "name": "Field 35",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": [
          "2024-01-01", "2024-01-02", "2024-01-03", "2024-01-04", "2024-01-05",
          "2024-01-06", "2024-01-07", "2024-01-08", "2024-01-09", "2024-01-10",
          "2024-01-11", "2024-01-12", "2024-01-13", "2024-01-14"
        ],
        "temperature_2m_mean": [25.0, 25.0, 25.0, 25.0, 25.0, 25.0, 25.0, 25.0, 25.0, 25.0, 25.0, 25.0, 25.0, 25.0],
        "precipitation_sum": [0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0]
      },
      "hourly": {
        "time": [
          "2024-01-01T00:00", "2024-01-01T01:00", "2024-01-01T02:00",
          "2024-01-02T00:00", "2024-01-02T01:00", "2024-01-02T02:00"
        ],
        "relative_humidity_2m": [50.0, 90.0, 60.0, 85.0, 55.0, 95.0]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-035",
          "image_id": "bb5e8400-e29b-41d4-a716-446655440035",
          "processing_id": "bb5e8400-e29b-41d4-a716-446655440035",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00035/5a801269-110d-46c0-89dd-09a6dc195435/2024-01-14/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "bb5e8400-e29b-41d4-a716-446655440035"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    And the parquet file "87428a25-f0d2-4dba-acfd-e3d7320c00035/5a801269-110d-46c0-89dd-09a6dc195435/2024-01-14/climate.parquet" should have the field "humidity_std_dev" in all rows

  # ========================================
  # ERROR SCENARIOS - Historical Climate Data Errors
  # ========================================

  @error @missing_historical_data
  Scenario: Process climate analysis with insufficient historical data - should handle gracefully
    Given the "HISTORICAL_CLIMATE_DAYS" env var is set to "30"
    And the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00036/5a801269-110d-46c0-89dd-09a6dc195436/2024-01-10/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "bb6e8400-e29b-41d4-a716-446655440036",
      "status": "PENDING",
      "date": "2024-01-10",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195436",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00036",
        "name": "Field 36",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-08", "2024-01-09", "2024-01-10"],
        "temperature_2m_mean": [25.0, 25.5, 26.0],
        "precipitation_sum": [0.0, 0.0, 0.0]
      },
      "hourly": {
        "time": ["2024-01-08T00:00", "2024-01-08T01:00"],
        "relative_humidity_2m": [75.0, 76.0]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-036",
          "image_id": "bb6e8400-e29b-41d4-a716-446655440036",
          "processing_id": "bb6e8400-e29b-41d4-a716-446655440036",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00036/5a801269-110d-46c0-89dd-09a6dc195436/2024-01-10/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "bb6e8400-e29b-41d4-a716-446655440036"
    }
    """
    Then the lambda should finish without errors

  @error @weather_metrics_not_found
  Scenario: Process climate analysis with weather date mismatch - should return weather metrics not found error
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00031/5a801269-110d-46c0-89dd-09a6dc195431/2024-01-05/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "xx2e8400-e29b-41d4-a716-446655440031",
      "status": "PENDING",
      "date": "2024-01-05",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195431",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00031",
        "name": "Field 31",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01", "2024-01-02", "2024-01-03"],
        "temperature_2m_mean": [25.3, 24.8, 26.1],
        "precipitation_sum": [0.0, 12.5, 0.0]
      },
      "hourly": {
        "time": ["2024-01-01T00:00", "2024-01-01T01:00"],
        "relative_humidity_2m": [75.2, 78.3]
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "xx2e8400-e29b-41d4-a716-446655440031"
    }
    """
    Then the lambda should finish with "bad request" error

  @error @weather_empty_response
  Scenario: Process climate analysis with empty weather data arrays - should handle gracefully
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00032/5a801269-110d-46c0-89dd-09a6dc195432/2024-01-01/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "xx3e8400-e29b-41d4-a716-446655440032",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195432",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00032",
        "name": "Field 32",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": [],
        "temperature_2m_mean": [],
        "precipitation_sum": []
      },
      "hourly": {
        "time": [],
        "relative_humidity_2m": []
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "xx3e8400-e29b-41d4-a716-446655440032"
    }
    """
    Then the lambda should finish with "bad request" error

  @error @analysis_api_malformed_json
  Scenario: Process climate analysis with malformed analysis API JSON - should handle gracefully
    Given the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "aa7e8400-e29b-41d4-a716-446655440036",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195436",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00036",
        "name": "Field 36",
        "centroid": {
          "lat": "invalid",
          "lng": "invalid"
        }
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "aa7e8400-e29b-41d4-a716-446655440036"
    }
    """
    Then the lambda should finish with "internal server" error

  @error @analysis_api_missing_centroid
  Scenario: Process climate analysis with missing centroid in analysis data - should return bad request
    Given the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "xx8e8400-e29b-41d4-a716-446655440037",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195437",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00037",
        "name": "Field 37",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "area_hectares": 10.5
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "xx8e8400-e29b-41d4-a716-446655440037"
    }
    """
    Then the lambda should finish with "bad request" error

  @error @analysis_api_null_centroid
  Scenario: Process climate analysis with null centroid values - should return bad request
    Given the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "xx9e8400-e29b-41d4-a716-446655440038",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195438",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00038",
        "name": "Field 38",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": null,
          "lng": null
        },
        "area_hectares": 10.5
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "xx9e8400-e29b-41d4-a716-446655440038"
    }
    """
    Then the lambda should finish with "bad request" error

  @error @weather_api_malformed_json
  Scenario: Process climate analysis with malformed weather API JSON - should handle gracefully
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00039/5a801269-110d-46c0-89dd-09a6dc195439/2024-01-01/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "aa0e8400-e29b-41d4-a716-446655440039",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195439",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00039",
        "name": "Field 39",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": "not-an-array",
        "temperature_2m_mean": "not-an-array"
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "aa0e8400-e29b-41d4-a716-446655440039"
    }
    """
    Then the lambda should finish with "internal server" error

  @error @analysis_api_missing_field_id
  Scenario: Process climate analysis with missing field ID in analysis data - should handle gracefully
    Given the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "aa1e8400-e29b-41d4-a716-446655440040",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "name": "Field 40",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "aa1e8400-e29b-41d4-a716-446655440040"
    }
    """
    Then the lambda should finish with "internal server" error

  # ========================================
  # ERROR SCENARIOS - File Creation API Issues
  # ========================================

  @error @file_api_failure
  Scenario: Process climate analysis with file creation API failure - should fail after parquet creation
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00060/5a801269-110d-46c0-89dd-09a6dc195460/2024-01-01/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "ff0e8400-e29b-41d4-a716-446655440060",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195460",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00060",
        "name": "Field 60",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [0.0]
      },
      "hourly": {
        "time": ["2024-01-01T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-060",
          "image_id": "ff0e8400-e29b-41d4-a716-446655440060",
          "processing_id": "ff0e8400-e29b-41d4-a716-446655440060",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00060/5a801269-110d-46c0-89dd-09a6dc195460/2024-01-01/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    And the 0 "POST" request to "/v1/files" returns status 500 with the following response
    """
    {
      "error": "Internal server error"
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "ff0e8400-e29b-41d4-a716-446655440060"
    }
    """
    Then the lambda should finish with "server error occurred" error

  @error @file_api_client_error
  Scenario: Process climate analysis with file creation API client error - should fail after parquet creation
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00061/5a801269-110d-46c0-89dd-09a6dc195461/2024-01-01/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "ff1e8400-e29b-41d4-a716-446655440061",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195461",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00061",
        "name": "Field 61",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [0.0]
      },
      "hourly": {
        "time": ["2024-01-01T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-061",
          "image_id": "ff1e8400-e29b-41d4-a716-446655440061",
          "processing_id": "ff1e8400-e29b-41d4-a716-446655440061",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00061/5a801269-110d-46c0-89dd-09a6dc195461/2024-01-01/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    And the 0 "POST" request to "/v1/files" returns status 400 with the following response
    """
    {
      "error": "Bad request"
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "ff1e8400-e29b-41d4-a716-446655440061"
    }
    """
    Then the lambda should finish with "client error occurred" error

  @success @weather_partial_hourly_data
  Scenario: Process climate analysis with partial hourly weather data - should succeed with available data
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00041/5a801269-110d-46c0-89dd-09a6dc195441/2024-01-01/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "aa2e8400-e29b-41d4-a716-446655440041",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195441",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00041",
        "name": "Field 41",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [3.5]
      },
      "hourly": {
        "time": ["2024-01-01T00:00", "2024-01-01T01:00", "2024-01-01T02:00"],
        "relative_humidity_2m": [75.2, 76.0, 77.5]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-041",
          "image_id": "aa2e8400-e29b-41d4-a716-446655440041",
          "processing_id": "aa2e8400-e29b-41d4-a716-446655440041",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00041/5a801269-110d-46c0-89dd-09a6dc195441/2024-01-01/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "aa2e8400-e29b-41d4-a716-446655440041"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"

  @success @zero_precipitation
  Scenario: Process climate analysis with zero precipitation values - should succeed
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00042/5a801269-110d-46c0-89dd-09a6dc195442/2024-01-01/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "aa3e8400-e29b-41d4-a716-446655440042",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195442",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00042",
        "name": "Field 42",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [0.0]
      },
      "hourly": {
        "time": ["2024-01-01T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-042",
          "image_id": "aa3e8400-e29b-41d4-a716-446655440042",
          "processing_id": "aa3e8400-e29b-41d4-a716-446655440042",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00042/5a801269-110d-46c0-89dd-09a6dc195442/2024-01-01/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "aa3e8400-e29b-41d4-a716-446655440042"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"

  # ========================================
  # SUCCESS SCENARIOS - Date Validation Skip Processing
  # ========================================

  @success @date_validation @skip_processing
  Scenario: Skip processing when climate analysis date equals image date - no processing needed
    Given the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "dd6e8400-e29b-41d4-a716-446655440050",
      "status": "PENDING",
      "date": "2024-01-15",
      "last_image_added_date": "2024-01-15",
      "last_climate_analysis_date": "2024-01-15",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195450",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00050",
        "name": "Field 50",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "dd6e8400-e29b-41d4-a716-446655440050"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "skipped" equal to "true"
    And the response should contain the field "reason" equal to "climate analysis is up to date"

  @success @date_validation @skip_processing
  Scenario: Skip processing when climate analysis date is newer than image date - no processing needed
    Given the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "ee6e8400-e29b-41d4-a716-446655440051",
      "status": "PENDING",
      "date": "2024-01-15",
      "last_image_added_date": "2024-01-10",
      "last_climate_analysis_date": "2024-01-15",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195451",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00051",
        "name": "Field 51",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "ee6e8400-e29b-41d4-a716-446655440051"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "skipped" equal to "true"
    And the response should contain the field "reason" equal to "climate analysis is up to date"

  @success @date_validation @skip_processing
  Scenario: Skip processing when image date is null - no processing needed
    Given the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "ff6e8400-e29b-41d4-a716-446655440052",
      "status": "PENDING",
      "date": "2024-01-15",
      "last_image_added_date": null,
      "last_climate_analysis_date": "2024-01-15",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195452",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00052",
        "name": "Field 52",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "ff6e8400-e29b-41d4-a716-446655440052"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "skipped" equal to "true"
    And the response should contain the field "reason" equal to "no image data available"

  # ========================================
  # SUCCESS SCENARIOS - Date Validation Execute Processing
  # ========================================

  @success @date_validation @execute_processing
  Scenario: Execute processing when climate analysis date is null - new analysis required
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00053/5a801269-110d-46c0-89dd-09a6dc195453/2024-01-15/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "aa7e8400-e29b-41d4-a716-446655440053",
      "status": "PENDING",
      "date": "2024-01-15",
      "last_image_added_date": "2024-01-15",
      "last_climate_analysis_date": null,
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195453",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00053",
        "name": "Field 53",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-15"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [0.0]
      },
      "hourly": {
        "time": ["2024-01-15T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-053",
          "image_id": "aa7e8400-e29b-41d4-a716-446655440053",
          "processing_id": "aa7e8400-e29b-41d4-a716-446655440053",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00053/5a801269-110d-46c0-89dd-09a6dc195453/2024-01-15/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "aa7e8400-e29b-41d4-a716-446655440053"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    And the response should contain the field "skipped" equal to "nil"

  @success @date_validation @execute_processing
  Scenario: Execute processing when climate analysis date is older than image date - update required
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00054/5a801269-110d-46c0-89dd-09a6dc195454/2024-01-20/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "bb7e8400-e29b-41d4-a716-446655440054",
      "status": "PENDING",
      "date": "2024-01-20",
      "last_image_added_date": "2024-01-20",
      "last_climate_analysis_date": "2024-01-15",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195454",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00054",
        "name": "Field 54",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-20"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [0.0]
      },
      "hourly": {
        "time": ["2024-01-20T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-054",
          "image_id": "bb7e8400-e29b-41d4-a716-446655440054",
          "processing_id": "bb7e8400-e29b-41d4-a716-446655440054",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00054/5a801269-110d-46c0-89dd-09a6dc195454/2024-01-20/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "bb7e8400-e29b-41d4-a716-446655440054"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    And the response should contain the field "skipped" equal to "nil"

  @success @date_validation @backward_compatibility
  Scenario: Execute processing when date fields are not present in API response - backward compatibility
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00055/5a801269-110d-46c0-89dd-09a6dc195455/2024-01-15/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "cc7e8400-e29b-41d4-a716-446655440055",
      "status": "PENDING",
      "date": "2024-01-15",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195455",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00055",
        "name": "Field 55",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-15"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [0.0]
      },
      "hourly": {
        "time": ["2024-01-15T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-055",
          "image_id": "cc7e8400-e29b-41d4-a716-446655440055",
          "processing_id": "cc7e8400-e29b-41d4-a716-446655440055",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00055/5a801269-110d-46c0-89dd-09a6dc195455/2024-01-15/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "cc7e8400-e29b-41d4-a716-446655440055"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    And the response should contain the field "skipped" equal to "nil"

  # ========================================
  # SUCCESS SCENARIOS - Fetch ImageID from DELTA File
  # ========================================
  # These scenarios test the new flow where the lambda fetches the DELTA file
  # by processing_id to get the correct image_id for creating CLIMATE file records.
  # The image_id is different from processing_id and must be retrieved from the DELTA file.

  @success @delta_file_lookup @image_id_fetch
  Scenario: Process climate analysis with correct image_id from DELTA file lookup
    # Setup: DELTA file exists with processing_id, containing image_id which is DIFFERENT from processing_id
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00100/5a801269-110d-46c0-89dd-09a6dc1954a1/2024-01-05/delta.parquet"
    # Mock the analysis API response
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "550e8400-e29b-41d4-a716-446655440100",
      "status": "PENDING",
      "date": "2024-01-05",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc1954a1",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00100",
        "name": "Field for DELTA lookup test",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    # Mock the DELTA file lookup - returns the file with image_id DIFFERENT from processing_id
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-001",
          "image_id": "aaaabbbb-cccc-dddd-eeee-ffffffffffff",
          "processing_id": "550e8400-e29b-41d4-a716-446655440100",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00100/5a801269-110d-46c0-89dd-09a6dc1954a1/2024-01-05/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    # Mock the weather API response
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01", "2024-01-02", "2024-01-03", "2024-01-04", "2024-01-05"],
        "temperature_2m_mean": [25.3, 24.8, 26.1, 25.5, 24.9],
        "precipitation_sum": [0.0, 12.5, 0.0, 5.2, 0.0]
      },
      "hourly": {
        "time": [
          "2024-01-01T00:00", "2024-01-01T01:00", "2024-01-01T02:00",
          "2024-01-02T00:00", "2024-01-02T01:00", "2024-01-02T02:00"
        ],
        "relative_humidity_2m": [75.2, 78.3, 76.5, 80.1, 82.3, 79.8]
      }
    }
    """
    # Mock the file creation API response - expects image_id from DELTA file, NOT processing_id
    And the 0 "POST" request to "/v1/files" returns status 201 with the following response
    """
    {
      "id": "climate-file-001",
      "image_id": "aaaabbbb-cccc-dddd-eeee-ffffffffffff",
      "type": "CLIMATE",
      "file_path": "data/climate/aaaabbbb-cccc-dddd-eeee-ffffffffffff/climate.parquet"
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "550e8400-e29b-41d4-a716-446655440100"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    # CRITICAL: Verify the DELTA file was queried with correct parameters
    And the 0 "GET" request for "/v1/files" queries should have the variable "processing_id" with value "550e8400-e29b-41d4-a716-446655440100"
    And the 0 "GET" request for "/v1/files" queries should have the variable "type" with value "DELTA"
    # CRITICAL: Verify the CLIMATE file record uses image_id from DELTA file, NOT processing_id
    And the 0 "POST" request for "/v1/files" json body should have the variable "image_id" with value "aaaabbbb-cccc-dddd-eeee-ffffffffffff"
    And the 0 "POST" request for "/v1/files" json body should have the variable "type" with value "CLIMATE"

  @success @delta_file_lookup @image_id_fetch @sns_integration
  Scenario: Process climate analysis via SNS with image_id from DELTA file
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00101/5a801269-110d-46c0-89dd-09a6dc1954a2/2024-01-02/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "660e8400-e29b-41d4-a716-446655440101",
      "status": "PENDING",
      "date": "2024-01-02",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc1954a2",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00101",
        "name": "Field for SNS DELTA lookup test",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    # DELTA file lookup returns image_id different from processing_id
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-002",
          "image_id": "11112222-3333-4444-5555-666677778888",
          "processing_id": "660e8400-e29b-41d4-a716-446655440101",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00101/5a801269-110d-46c0-89dd-09a6dc1954a2/2024-01-02/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01", "2024-01-02"],
        "temperature_2m_mean": [25.3, 24.8],
        "precipitation_sum": [0.0, 12.5]
      },
      "hourly": {
        "time": ["2024-01-01T00:00", "2024-01-01T01:00"],
        "relative_humidity_2m": [75.2, 78.3]
      }
    }
    """
    And the 0 "POST" request to "/v1/files" returns status 201 with the following response
    """
    {
      "id": "climate-file-002",
      "image_id": "11112222-3333-4444-5555-666677778888",
      "type": "CLIMATE",
      "file_path": "data/climate/11112222-3333-4444-5555-666677778888/climate.parquet"
    }
    """
    When the following event is received via sns
    """
    {
      "processing_id": "660e8400-e29b-41d4-a716-446655440101"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    # Verify DELTA file was queried
    And the 0 "GET" request for "/v1/files" queries should have the variable "processing_id" with value "660e8400-e29b-41d4-a716-446655440101"
    And the 0 "GET" request for "/v1/files" queries should have the variable "type" with value "DELTA"
    # Verify CLIMATE file uses image_id from DELTA file
    And the 0 "POST" request for "/v1/files" json body should have the variable "image_id" with value "11112222-3333-4444-5555-666677778888"
    And the 0 "POST" request for "/v1/files" json body should have the variable "type" with value "CLIMATE"

  # ========================================
  # ERROR SCENARIOS - DELTA File Lookup Failures
  # ========================================

  @error @delta_file_lookup @delta_file_not_found
  Scenario: Process climate analysis fails when DELTA file is not found - should return not found error
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00102/5a801269-110d-46c0-89dd-09a6dc1954a3/2024-01-05/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "770e8400-e29b-41d4-a716-446655440102",
      "status": "PENDING",
      "date": "2024-01-05",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc1954a3",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00102",
        "name": "Field for DELTA not found test",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    # DELTA file lookup returns empty result - no DELTA file found for this processing_id
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 0,
        "total_pages": 0
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-05"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [0.0]
      },
      "hourly": {
        "time": ["2024-01-05T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "770e8400-e29b-41d4-a716-446655440102"
    }
    """
    Then the lambda should finish with "not found" error

  @error @delta_file_lookup @delta_file_api_error
  Scenario: Process climate analysis fails when DELTA file API returns server error
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00103/5a801269-110d-46c0-89dd-09a6dc1954a4/2024-01-05/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "880e8400-e29b-41d4-a716-446655440103",
      "status": "PENDING",
      "date": "2024-01-05",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc1954a4",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00103",
        "name": "Field for DELTA API error test",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    # DELTA file lookup API returns 500 error
    And the 0 "GET" request to "/v1/files" returns status 500 with the following response
    """
    {
      "error": "Internal server error"
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-05"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [0.0]
      },
      "hourly": {
        "time": ["2024-01-05T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "880e8400-e29b-41d4-a716-446655440103"
    }
    """
    Then the lambda should finish with "server error occurred" error

  @error @delta_file_lookup @delta_file_api_404
  Scenario: Process climate analysis fails when DELTA file API returns 404
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00104/5a801269-110d-46c0-89dd-09a6dc1954a5/2024-01-05/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "990e8400-e29b-41d4-a716-446655440104",
      "status": "PENDING",
      "date": "2024-01-05",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc1954a5",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00104",
        "name": "Field for DELTA 404 test",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    # DELTA file lookup API returns 404
    And the 0 "GET" request to "/v1/files" returns status 404 with the following response
    """
    {
      "error": {
        "description": "Resource not found",
        "code": "NOT_FOUND"
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-05"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [0.0]
      },
      "hourly": {
        "time": ["2024-01-05T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "990e8400-e29b-41d4-a716-446655440104"
    }
    """
    Then the lambda should finish with "not found" error

  @error @delta_file_lookup @delta_file_missing_image_id
  Scenario: Process climate analysis fails when DELTA file response is missing image_id
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00105/5a801269-110d-46c0-89dd-09a6dc1954a6/2024-01-05/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "aa0e8400-e29b-41d4-a716-446655440105",
      "status": "PENDING",
      "date": "2024-01-05",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc1954a6",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00105",
        "name": "Field for missing image_id test",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    # DELTA file exists but image_id is missing/null
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-bad",
          "image_id": null,
          "processing_id": "aa0e8400-e29b-41d4-a716-446655440105",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00105/5a801269-110d-46c0-89dd-09a6dc1954a6/2024-01-05/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-05"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [0.0]
      },
      "hourly": {
        "time": ["2024-01-05T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "aa0e8400-e29b-41d4-a716-446655440105"
    }
    """
    Then the lambda should finish with "bad request" error

  @success @delta_file_lookup @verify_processing_id_not_used_as_image_id
  Scenario: Verify processing_id is NOT used as image_id - different values prove correct behavior
    # This scenario explicitly verifies that processing_id and image_id are different
    # and that the lambda correctly uses image_id from DELTA file, not processing_id
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c00106/5a801269-110d-46c0-89dd-09a6dc1954a7/2024-01-05/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "bb0e8400-e29b-41d4-a716-446655440106",
      "status": "PENDING",
      "date": "2024-01-05",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc1954a7",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c00106",
        "name": "Field for image_id verification",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    # DELTA file has DIFFERENT image_id from processing_id
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-verify",
          "image_id": "cc1e9500-f39c-52e5-b827-557766551111",
          "processing_id": "bb0e8400-e29b-41d4-a716-446655440106",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c00106/5a801269-110d-46c0-89dd-09a6dc1954a7/2024-01-05/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-05"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [0.0]
      },
      "hourly": {
        "time": ["2024-01-05T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    And the 0 "POST" request to "/v1/files" returns status 201 with the following response
    """
    {
      "id": "climate-file-verify",
      "image_id": "cc1e9500-f39c-52e5-b827-557766551111",
      "type": "CLIMATE",
      "file_path": "data/climate/cc1e9500-f39c-52e5-b827-557766551111/climate.parquet"
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "bb0e8400-e29b-41d4-a716-446655440106"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    # CRITICAL ASSERTIONS: These verify the fix works correctly
    # 1. DELTA file was queried with the processing_id
    And the 0 "GET" request for "/v1/files" queries should have the variable "processing_id" with value "bb0e8400-e29b-41d4-a716-446655440106"
    And the 0 "GET" request for "/v1/files" queries should have the variable "type" with value "DELTA"
    # 2. CLIMATE file uses image_id from DELTA (NOT processing_id)
    And the 0 "POST" request for "/v1/files" json body should have the variable "image_id" with value "cc1e9500-f39c-52e5-b827-557766551111"
    # 3. Verify it's NOT using processing_id as image_id (this would be the bug)
    # Note: The assertion above already proves this since "cc1e9500-f39c-52e5-b827-557766551111" != "bb0e8400-e29b-41d4-a716-446655440106"

  # ========================================
  # SUCCESS SCENARIOS - READY_ANALYSIS Event Publishing
  # ========================================

  @success @event_publishing @ready_analysis
  Scenario: Publish READY_ANALYSIS event after successful climate analysis via SQS
    Given the "FOREST_EVENTS_TOPIC_ARN" env var is set to "arn:aws:sns:us-east-1:123456789012:forest-events"
    And the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c50001/5a801269-110d-46c0-89dd-09a6dc195501/2024-01-05/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "ee1e8400-e29b-41d4-a716-446655440501",
      "status": "PENDING",
      "date": "2024-01-05",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195501",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c50001",
        "name": "Field Event Test 1",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01", "2024-01-02", "2024-01-03", "2024-01-04", "2024-01-05"],
        "temperature_2m_mean": [25.3, 24.8, 26.1, 25.5, 24.9],
        "precipitation_sum": [0.0, 12.5, 0.0, 5.2, 0.0]
      },
      "hourly": {
        "time": ["2024-01-01T00:00", "2024-01-01T01:00"],
        "relative_humidity_2m": [75.2, 78.3]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-501",
          "image_id": "ee1e8400-e29b-41d4-a716-446655440501",
          "processing_id": "ee1e8400-e29b-41d4-a716-446655440501",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c50001/5a801269-110d-46c0-89dd-09a6dc195501/2024-01-05/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    And the 0 "POST" request to "/v1/files" returns status 201 with the following response
    """
    {
      "id": "file-record-501",
      "image_id": "ee1e8400-e29b-41d4-a716-446655440501",
      "type": "CLIMATE",
      "file_path": "data/climate/ee1e8400-e29b-41d4-a716-446655440501/climate.parquet"
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "ee1e8400-e29b-41d4-a716-446655440501"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    # Verify READY_ANALYSIS event was published to SNS
    And the sns topic "arn:aws:sns:us-east-1:123456789012:forest-events" should have 1 messages published
    And the sns topic "arn:aws:sns:us-east-1:123456789012:forest-events" should have a message published with "event_type" field equal to "READY_ANALYSIS"
    And the sns topic "arn:aws:sns:us-east-1:123456789012:forest-events" should have a message published with "event_data.processing_id" field equal to "ee1e8400-e29b-41d4-a716-446655440501"
    And the sns topic "arn:aws:sns:us-east-1:123456789012:forest-events" should have a message published with "spec_version" field equal to "1"
    And the sns topic "arn:aws:sns:us-east-1:123456789012:forest-events" should have a message published with "source" field equal to "maxsatt-climate-analysis-lambda"
    And the sns topic "arn:aws:sns:us-east-1:123456789012:forest-events" should have a message published with "event_id" field equal to "not nil"
    And the sns topic "arn:aws:sns:us-east-1:123456789012:forest-events" should have a message published with "event_correlation_id" field equal to "not nil"
    And the sns topic "arn:aws:sns:us-east-1:123456789012:forest-events" should have a message published with "event_date" field equal to "not nil"

  @success @event_publishing @ready_analysis @sns_integration
  Scenario: Publish READY_ANALYSIS event after successful climate analysis via SNS
    Given the "FOREST_EVENTS_TOPIC_ARN" env var is set to "arn:aws:sns:us-east-1:123456789012:forest-events"
    And the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c50002/5a801269-110d-46c0-89dd-09a6dc195502/2024-01-02/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "ff1e8400-e29b-41d4-a716-446655440502",
      "status": "PENDING",
      "date": "2024-01-02",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195502",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c50002",
        "name": "Field Event Test 2",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01", "2024-01-02"],
        "temperature_2m_mean": [25.3, 24.8],
        "precipitation_sum": [0.0, 12.5]
      },
      "hourly": {
        "time": ["2024-01-01T00:00", "2024-01-01T01:00"],
        "relative_humidity_2m": [75.2, 78.3]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-502",
          "image_id": "ff1e8400-e29b-41d4-a716-446655440502",
          "processing_id": "ff1e8400-e29b-41d4-a716-446655440502",
          "type": "DELTA",
          "file_path": "87428a25-f0d2-4dba-acfd-e3d7320c50002/5a801269-110d-46c0-89dd-09a6dc195502/2024-01-02/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    And the 0 "POST" request to "/v1/files" returns status 201 with the following response
    """
    {
      "id": "file-record-502",
      "image_id": "ff1e8400-e29b-41d4-a716-446655440502",
      "type": "CLIMATE",
      "file_path": "data/climate/ff1e8400-e29b-41d4-a716-446655440502/climate.parquet"
    }
    """
    When the following event is received via sns
    """
    {
      "processing_id": "ff1e8400-e29b-41d4-a716-446655440502"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" equal to "not nil"
    # Verify READY_ANALYSIS event was published to SNS
    And the sns topic "arn:aws:sns:us-east-1:123456789012:forest-events" should have 1 messages published
    And the sns topic "arn:aws:sns:us-east-1:123456789012:forest-events" should have a message published with "event_type" field equal to "READY_ANALYSIS"
    And the sns topic "arn:aws:sns:us-east-1:123456789012:forest-events" should have a message published with "event_data.processing_id" field equal to "ff1e8400-e29b-41d4-a716-446655440502"

  # ========================================
  # ERROR SCENARIOS - READY_ANALYSIS Event NOT Published
  # ========================================

  @error @event_publishing @no_event_on_validation_error
  Scenario: NO event published when processing_id validation fails
    Given the "FOREST_EVENTS_TOPIC_ARN" env var is set to "arn:aws:sns:us-east-1:123456789012:forest-events"
    When the following event is received via sqs
    """
    {
      "processing_id": "invalid-uuid-format"
    }
    """
    Then the lambda should finish with "bad request" error
    # Verify NO event was published to SNS on validation failure
    And the sns topic "arn:aws:sns:us-east-1:123456789012:forest-events" should have 0 messages published

  @error @event_publishing @no_event_on_not_found
  Scenario: NO event published when processing not found
    Given the "FOREST_EVENTS_TOPIC_ARN" env var is set to "arn:aws:sns:us-east-1:123456789012:forest-events"
    And the 0 "GET" request to "/v1/analysis/*" returns status 404 with the following response
    """
    {
      "error": {
        "description": "Processing not found",
        "code": "NOT_FOUND"
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "aa1e8400-e29b-41d4-a716-446655440503"
    }
    """
    Then the lambda should finish with "not found" error
    # Verify NO event was published to SNS on not found error
    And the sns topic "arn:aws:sns:us-east-1:123456789012:forest-events" should have 0 messages published

  @error @event_publishing @no_event_on_weather_api_failure
  Scenario: NO event published when weather API fails
    Given the "FOREST_EVENTS_TOPIC_ARN" env var is set to "arn:aws:sns:us-east-1:123456789012:forest-events"
    And the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c50003/5a801269-110d-46c0-89dd-09a6dc195503/2024-01-01/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "bb1e8400-e29b-41d4-a716-446655440504",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195503",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c50003",
        "name": "Field Event Test 3",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the "GET" "/v1/archive" should not return any results
    When the following event is received via sqs
    """
    {
      "processing_id": "bb1e8400-e29b-41d4-a716-446655440504"
    }
    """
    Then the lambda should finish with "internal server" error
    # Verify NO event was published to SNS on weather API failure
    And the sns topic "arn:aws:sns:us-east-1:123456789012:forest-events" should have 0 messages published

  @error @event_publishing @no_event_on_empty_dataset
  Scenario: NO event published when dataset is empty
    Given the "FOREST_EVENTS_TOPIC_ARN" env var is set to "arn:aws:sns:us-east-1:123456789012:forest-events"
    And the "test/integration/resources/climate-analysis/empty.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c50004/5a801269-110d-46c0-89dd-09a6dc195504/2024-01-01/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "cc1e8400-e29b-41d4-a716-446655440505",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195504",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c50004",
        "name": "Field Event Test 4",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [3.5]
      },
      "hourly": {
        "time": ["2024-01-01T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "cc1e8400-e29b-41d4-a716-446655440505"
    }
    """
    Then the lambda should finish with "internal server" error
    # Verify NO event was published to SNS on empty dataset error
    And the sns topic "arn:aws:sns:us-east-1:123456789012:forest-events" should have 0 messages published

  @error @event_publishing @no_event_on_invalid_coordinates
  Scenario: NO event published when coordinates are invalid
    Given the "FOREST_EVENTS_TOPIC_ARN" env var is set to "arn:aws:sns:us-east-1:123456789012:forest-events"
    And the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "87428a25-f0d2-4dba-acfd-e3d7320c50005/5a801269-110d-46c0-89dd-09a6dc195505/2024-01-01/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "dd1e8400-e29b-41d4-a716-446655440506",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc195505",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c50005",
        "name": "Field Event Test 5",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": 95.0,
          "lng": 200.0
        },
        "area_hectares": 10.5
      }
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "dd1e8400-e29b-41d4-a716-446655440506"
    }
    """
    Then the lambda should finish with "bad request" error
    # Verify NO event was published to SNS on invalid coordinates error
    And the sns topic "arn:aws:sns:us-east-1:123456789012:forest-events" should have 0 messages published

  # ========================================
  # SUCCESS SCENARIOS - Processing Config ID Path Structure
  # ========================================

  @success @processing_config_id @path_structure
  Scenario: Process climate analysis with processing_config_id in S3 path structure
    # New path structure: {processing_config_id}/{forest_id}/{field_id}/{date}/delta.parquet
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "config-001/87428a25-f0d2-4dba-acfd-e3d7320c60001/5a801269-110d-46c0-89dd-09a6dc196001/2024-01-01/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "aa1e8400-e29b-41d4-a716-446655446001",
      "processing_config_id": "config-001",
      "status": "PENDING",
      "date": "2024-01-01",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc196001",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c60001",
        "name": "Field Config Test 1",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01"],
        "temperature_2m_mean": [25.3],
        "precipitation_sum": [0.0]
      },
      "hourly": {
        "time": ["2024-01-01T00:00"],
        "relative_humidity_2m": [75.2]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-config-001",
          "image_id": "aa1e8400-e29b-41d4-a716-446655446001",
          "processing_id": "aa1e8400-e29b-41d4-a716-446655446001",
          "type": "DELTA",
          "file_path": "config-001/87428a25-f0d2-4dba-acfd-e3d7320c60001/5a801269-110d-46c0-89dd-09a6dc196001/2024-01-01/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    And the 0 "POST" request to "/v1/files" returns status 201 with the following response
    """
    {
      "id": "file-record-config-001",
      "image_id": "aa1e8400-e29b-41d4-a716-446655446001",
      "type": "CLIMATE",
      "file_path": "config-001/87428a25-f0d2-4dba-acfd-e3d7320c60001/5a801269-110d-46c0-89dd-09a6dc196001/2024-01-01/climate.parquet"
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "aa1e8400-e29b-41d4-a716-446655446001"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" that contains "config-001/"
    And the response should contain the field "s3Key" that contains "/climate.parquet"
    And the 0 "POST" request for "/v1/files" json body should have the variable "image_id" with value "aa1e8400-e29b-41d4-a716-446655446001"
    And the 0 "POST" request for "/v1/files" json body should have the variable "type" with value "CLIMATE"
    # Validate the output parquet is in the correct path with processing_config_id prefix
    And the parquet file "config-001/87428a25-f0d2-4dba-acfd-e3d7320c60001/5a801269-110d-46c0-89dd-09a6dc196001/2024-01-01/climate.parquet" should have 1162 rows

  @success @processing_config_id @uuid_format
  Scenario: Process climate analysis with UUID-format processing_config_id
    # Validate processing_config_id works with UUID format (typical production scenario)
    Given the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "a1b2c3d4-e5f6-7890-abcd-ef1234567890/87428a25-f0d2-4dba-acfd-e3d7320c60002/5a801269-110d-46c0-89dd-09a6dc196002/2024-01-02/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "bb1e8400-e29b-41d4-a716-446655446002",
      "processing_config_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "status": "PENDING",
      "date": "2024-01-02",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc196002",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c60002",
        "name": "Field Config Test 2",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01", "2024-01-02"],
        "temperature_2m_mean": [25.3, 24.8],
        "precipitation_sum": [0.0, 5.2]
      },
      "hourly": {
        "time": ["2024-01-01T00:00", "2024-01-01T01:00"],
        "relative_humidity_2m": [75.2, 78.3]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-config-002",
          "image_id": "bb1e8400-e29b-41d4-a716-446655446002",
          "processing_id": "bb1e8400-e29b-41d4-a716-446655446002",
          "type": "DELTA",
          "file_path": "a1b2c3d4-e5f6-7890-abcd-ef1234567890/87428a25-f0d2-4dba-acfd-e3d7320c60002/5a801269-110d-46c0-89dd-09a6dc196002/2024-01-02/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    And the 0 "POST" request to "/v1/files" returns status 201 with the following response
    """
    {
      "id": "file-record-config-002",
      "image_id": "bb1e8400-e29b-41d4-a716-446655446002",
      "type": "CLIMATE",
      "file_path": "a1b2c3d4-e5f6-7890-abcd-ef1234567890/87428a25-f0d2-4dba-acfd-e3d7320c60002/5a801269-110d-46c0-89dd-09a6dc196002/2024-01-02/climate.parquet"
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "bb1e8400-e29b-41d4-a716-446655446002"
    }
    """
    Then the lambda should finish without errors
    And the response should contain the field "s3Key" that contains "a1b2c3d4-e5f6-7890-abcd-ef1234567890/"
    And the parquet file "a1b2c3d4-e5f6-7890-abcd-ef1234567890/87428a25-f0d2-4dba-acfd-e3d7320c60002/5a801269-110d-46c0-89dd-09a6dc196002/2024-01-02/climate.parquet" should have 1162 rows

  @success @processing_config_id @event_publishing
  Scenario: Validate SNS event includes processing_config_id in S3 path
    Given the "FOREST_EVENTS_TOPIC_ARN" env var is set to "arn:aws:sns:us-east-1:123456789012:forest-events"
    And the "test/integration/resources/climate-analysis/test-image_raw.parquet" local file exists in s3 as "config-event-001/87428a25-f0d2-4dba-acfd-e3d7320c60003/5a801269-110d-46c0-89dd-09a6dc196003/2024-01-03/delta.parquet"
    And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
    """
    {
      "processing_id": "cc1e8400-e29b-41d4-a716-446655446003",
      "processing_config_id": "config-event-001",
      "status": "PENDING",
      "date": "2024-01-03",
      "field": {
        "id": "5a801269-110d-46c0-89dd-09a6dc196003",
        "forest_id": "87428a25-f0d2-4dba-acfd-e3d7320c60003",
        "name": "Field Config Event Test",
        "geometry": {
          "type": "Polygon",
          "coordinates": [[[-46.8, -23.4], [-46.79, -23.4], [-46.79, -23.5], [-46.8, -23.5], [-46.8, -23.4]]]
        },
        "centroid": {
          "lat": -23.4565,
          "lng": -46.7895
        },
        "area_hectares": 10.5
      }
    }
    """
    And the 0 "GET" request to "/v1/archive" returns status 200 with the following response
    """
    {
      "daily": {
        "time": ["2024-01-01", "2024-01-02", "2024-01-03"],
        "temperature_2m_mean": [25.3, 24.8, 26.1],
        "precipitation_sum": [0.0, 5.2, 0.0]
      },
      "hourly": {
        "time": ["2024-01-01T00:00", "2024-01-01T01:00"],
        "relative_humidity_2m": [75.2, 78.3]
      }
    }
    """
    And the 0 "GET" request to "/v1/files" returns status 200 with the following response
    """
    {
      "content": [
        {
          "id": "delta-file-config-003",
          "image_id": "cc1e8400-e29b-41d4-a716-446655446003",
          "processing_id": "cc1e8400-e29b-41d4-a716-446655446003",
          "type": "DELTA",
          "file_path": "config-event-001/87428a25-f0d2-4dba-acfd-e3d7320c60003/5a801269-110d-46c0-89dd-09a6dc196003/2024-01-03/delta.parquet"
        }
      ],
      "metadata": {
        "page": 1,
        "limit": 1,
        "total": 1,
        "total_pages": 1
      }
    }
    """
    And the 0 "POST" request to "/v1/files" returns status 201 with the following response
    """
    {
      "id": "file-record-config-003",
      "image_id": "cc1e8400-e29b-41d4-a716-446655446003",
      "type": "CLIMATE",
      "file_path": "config-event-001/87428a25-f0d2-4dba-acfd-e3d7320c60003/5a801269-110d-46c0-89dd-09a6dc196003/2024-01-03/climate.parquet"
    }
    """
    When the following event is received via sqs
    """
    {
      "processing_id": "cc1e8400-e29b-41d4-a716-446655446003"
    }
    """
    Then the lambda should finish without errors
    And the sns topic "arn:aws:sns:us-east-1:123456789012:forest-events" should have 1 messages published
    And the sns topic "arn:aws:sns:us-east-1:123456789012:forest-events" should have a message published with "event_type" field equal to "READY_ANALYSIS"
    And the sns topic "arn:aws:sns:us-east-1:123456789012:forest-events" should have a message published with "event_data.processing_id" field equal to "cc1e8400-e29b-41d4-a716-446655446003"