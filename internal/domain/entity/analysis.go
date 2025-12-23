package entity

import "time"

// Analysis represents the analysis data retrieved from the external API.
// It contains field information with geometry and pre-calculated centroid.
type Analysis struct {
	Date                    time.Time
	LastImageAddedDate      *time.Time
	LastClimateAnalysisDate *time.Time
	Geometry                *Geometry
	Centroid                *Centroid
	ProcessingID            string
	ProcessingConfigID      string
	Status                  string
	ForestID                string
	FieldID                 string
	FieldName               string
	AreaHectares            float64
}

// HasCentroid checks if the analysis data includes a valid centroid.
func (a *Analysis) HasCentroid() bool {
	return a.Centroid != nil && a.Centroid.IsValid()
}

// IsCentroidValid checks if the centroid is valid.
func (a *Analysis) IsCentroidValid() bool {
	return a.HasCentroid()
}

// Longitude returns the longitude of the centroid if available, otherwise returns 0.0.
func (a *Analysis) Longitude() float64 {
	if a.HasCentroid() {
		return a.Centroid.Longitude
	}
	return 0.0
}

// Latitude returns the latitude of the centroid if available, otherwise returns 0.0.
func (a *Analysis) Latitude() float64 {
	if a.HasCentroid() {
		return a.Centroid.Latitude
	}
	return 0.0
}

// ShouldSkipProcessing determines if climate analysis processing should be skipped.
// Returns true with a reason if processing should be skipped, false with empty reason otherwise.
// Business rules:
// 1. Backward compatibility: If both date fields are null, execute processing (original behavior)
// 2. Skip if last_image_added_date is null but last_climate_analysis_date exists - no image data available
// 3. Skip if last_climate_analysis_date >= last_image_added_date - analysis is up to date
// 4. Execute if last_climate_analysis_date is null or older than last_image_added_date
func (a *Analysis) ShouldSkipProcessing() (skip bool, reason string) {
	if a.LastImageAddedDate == nil && a.LastClimateAnalysisDate == nil {
		return false, ""
	}

	if a.LastImageAddedDate == nil {
		return true, "no image data available"
	}

	if a.LastClimateAnalysisDate != nil {
		if !a.LastClimateAnalysisDate.Before(*a.LastImageAddedDate) {
			return true, "climate analysis is up to date"
		}
	}

	return false, ""
}

// DatePath returns the date formatted for use in S3 path (yyyy-MM-dd format).
func (a *Analysis) DatePath() string {
	return a.Date.Format("2006-01-02")
}

// BasePath returns the base S3 path for this analysis.
// Format varies based on ProcessingConfigID availability:
//   - With ProcessingConfigID: {processing_config_id}/{forest_id}/{field_id}/{date(yyyy-MM-dd)}
//   - Without ProcessingConfigID: {forest_id}/{field_id}/{date(yyyy-MM-dd)} (backward compatibility)
func (a *Analysis) BasePath() string {
	if a.ProcessingConfigID == "" {
		return a.ForestID + "/" + a.FieldID + "/" + a.DatePath()
	}
	return a.ProcessingConfigID + "/" + a.ForestID + "/" + a.FieldID + "/" + a.DatePath()
}
