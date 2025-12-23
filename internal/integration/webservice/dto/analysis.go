package dto

import (
	"time"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/timeutils"
)

// Analysis represents the API response when fetching analysis data by processing_id.
// It matches the structure returned by the external analysis API.
type Analysis struct {
	Date                    timeutils.FlexDate  `json:"date"`
	LastImageAddedDate      *timeutils.FlexDate `json:"last_image_added_date"`
	LastClimateAnalysisDate *timeutils.FlexDate `json:"last_climate_analysis_date"`
	CreatedAt               time.Time           `json:"created_at"`
	UpdatedAt               time.Time           `json:"updated_at"`
	ProcessingID            string              `json:"processing_id"`
	ProcessingConfigID      *string             `json:"processing_config_id"`
	Status                  string              `json:"status"`
	Field                   FieldData           `json:"field"`
}

// Geometry represents a GeoJSON geometry object.
// It contains the type and coordinates of a geographic feature.
// Coordinates is any type to support different GeoJSON geometry types (Point, Polygon, etc).
type Geometry struct {
	Coordinates any    `json:"coordinates"`
	Type        string `json:"type"`
}

// FieldData represents the field information within the analysis response.
type FieldData struct {
	Geometry     Geometry     `json:"geometry"`
	ID           string       `json:"id"`
	ForestID     string       `json:"forest_id"`
	Name         string       `json:"name"`
	Centroid     CentroidData `json:"centroid"`
	AreaHectares float64      `json:"area_hectares"`
}

// CentroidData represents the centroid coordinates from the API.
type CentroidData struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

func (a *Analysis) ToEntity() entity.Analysis {
	var coordinates []float64
	if coordsAny, ok := a.Field.Geometry.Coordinates.([]any); ok && len(coordsAny) > 0 {
		coordinates = nil
	}

	var lastImageAddedDate *time.Time
	if a.LastImageAddedDate != nil {
		lastImageAddedDate = &a.LastImageAddedDate.Time
	}

	var lastClimateAnalysisDate *time.Time
	if a.LastClimateAnalysisDate != nil {
		lastClimateAnalysisDate = &a.LastClimateAnalysisDate.Time
	}

	var processingConfigID string
	if a.ProcessingConfigID != nil {
		processingConfigID = *a.ProcessingConfigID
	}

	return entity.Analysis{
		ProcessingID:            a.ProcessingID,
		ProcessingConfigID:      processingConfigID,
		Status:                  a.Status,
		Date:                    a.Date.Time,
		LastImageAddedDate:      lastImageAddedDate,
		LastClimateAnalysisDate: lastClimateAnalysisDate,
		ForestID:                a.Field.ForestID,
		FieldID:                 a.Field.ID,
		FieldName:               a.Field.Name,
		Geometry: &entity.Geometry{
			Type:        a.Field.Geometry.Type,
			Coordinates: coordinates,
		},
		Centroid: &entity.Centroid{
			Latitude:  a.Field.Centroid.Latitude,
			Longitude: a.Field.Centroid.Longitude,
		},
		AreaHectares: a.Field.AreaHectares,
	}
}
