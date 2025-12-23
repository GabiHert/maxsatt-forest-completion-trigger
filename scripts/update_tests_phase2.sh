#!/bin/bash

# Script to update all BDD test scenarios for Phase 2 refactoring
# This script updates the feature file to:
# 1. Replace JSON file references with parquet files
# 2. Add Analysis API mocks to all scenarios
# 3. Update field assertions from forest/plot to field_name/field_id

FEATURE_FILE="test/integration/features/climate-analysis.feature"

echo "Updating BDD test scenarios for Phase 2..."

# Backup the original file
cp "$FEATURE_FILE" "$FEATURE_FILE.backup"

# Use sed to make systematic replacements
# Note: These are complex multiline replacements, so we'll use perl instead of sed

perl -i -pe '
  # Replace delta_dataset.json with parquet file
  s|delta_dataset_valid\.json.*as "data/processing/([^/]+)/delta_dataset\.json"|test-image_raw.parquet" local file exists in s3 as "data/processing/$1/test-image_raw.parquet"|g;

  # Replace empty delta dataset
  s|delta_dataset_empty\.json.*as "data/processing/([^/]+)/delta_dataset\.json"|test-image_raw.parquet" local file exists in s3 as "data/processing/$1/test-image_raw.parquet"|g;

  # Replace malformed delta dataset
  s|delta_dataset_malformed\.json.*as "data/processing/([^/]+)/delta_dataset\.json"|test-image_raw.parquet" local file exists in s3 as "data/processing/$1/test-image_raw.parquet"|g;

  # Remove geometry file lines (no longer needed)
  s|^\s*And the "test/integration/resources/climate-analysis/geometry_.*\.geojson".*\n||g;

  # Remove metadata file lines (no longer needed)
  s|^\s*And the "test/integration/resources/climate-analysis/metadata_.*\.json".*\n||g;

  # Replace forest assertions with field_name
  s|finalData\.0\.forest.*equal to "forest_01"|finalData.0.field_name" equal to "Field 1"|g;

  # Replace plot assertions with field_id
  s|finalData\.0\.plot.*equal to "plot_12"|finalData.0.field_id" equal to "field-001"|g;

' "$FEATURE_FILE"

echo "Phase 1 replacements complete. Now adding Analysis API mocks..."

# The Analysis API mocks need to be added before the weather API mock in each scenario
# This is complex, so we'll document it for manual addition

cat > "test/integration/PHASE2_UPDATES_NEEDED.md" << 'EOF'
# Phase 2 Test Updates Still Needed

## Remaining Manual Updates Required

### 1. Add Analysis API Mocks to All Scenarios

Each scenario that tests the success path needs an Analysis API mock BEFORE the weather API mock.
Add this block after the S3 file setup and before the weather API mock:

```
And the 0 "GET" request to "/v1/analysis/*" returns status 200 with the following response
"""
{
  "processing_id": "<PROCESSING_ID>",
  "status": "PENDING",
  "field": {
    "id": "field-001",
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
```

Replace `<PROCESSING_ID>` with the actual processing ID from each scenario.

### 2. Scenarios That Need Analysis API Mocks

All scenarios in these sections:
- SUCCESS SCENARIOS - Complete Pipeline
- SUCCESS SCENARIOS - Centroid & Weather API
- SUCCESS SCENARIOS - Response Structure
- SUCCESS SCENARIOS - Weather Metrics
- SUCCESS SCENARIOS - Caching
- SUCCESS SCENARIOS - Concurrent Processing

### 3. Error Scenarios

Error scenarios that test "not found" or "API failures" may need adjustments:
- Missing S3 files tests → Now test for missing parquet
- Geometry-related errors → Now come from Analysis API response

### 4. Field Name Assertions

Update any remaining assertions:
- `forest` → `field_name`
- `plot` → `field_id`

Expected values:
- `forest_01` → `Field 1`
- `plot_12` → `field-001`

EOF

echo ""
echo "Script complete!"
echo "Basic replacements have been made."
echo "See test/integration/PHASE2_UPDATES_NEEDED.md for remaining manual updates."
echo ""
echo "Backup saved to: $FEATURE_FILE.backup"
