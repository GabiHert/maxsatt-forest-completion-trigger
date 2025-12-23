package entity

// Geometry represents a GeoJSON geometry object.
// It contains the type and coordinates of a geographic feature.
type Geometry struct {
	Type        string
	Coordinates []float64
}

// Centroid represents a geographic center point with latitude and longitude.
type Centroid struct {
	Latitude  float64
	Longitude float64
}

// IsValid checks if the centroid coordinates are within valid ranges.
// Latitude must be between -90 and 90, longitude between -180 and 180.
func (c *Centroid) IsValid() bool {
	return c.Latitude >= -90 && c.Latitude <= 90 &&
		c.Longitude >= -180 && c.Longitude <= 180
}
