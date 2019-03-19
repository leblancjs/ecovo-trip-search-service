package entity

import (
	"fmt"
)

// Point contains a geolocation's information.
type Point struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Name      string  `json:"name,ommitempty"`
}

const (
	// MinimumLongitude represents the minimum longitude value.
	MinimumLongitude = -180

	// MaximumLongitude represents the maximum longitude value.
	MaximumLongitude = 180

	// MinimumLatitude represents the minimum latitude value.
	MinimumLatitude = -90

	// MaximumLatitude represents the maximum latitude value.
	MaximumLatitude = 90
)

// String returns string value of Point.
func (m *Point) String() string {
	return fmt.Sprintf("%f", m.Latitude) + ", " + fmt.Sprintf("%f", m.Longitude)
}

// Validate validates that the map's required fields are filled out correctly.
func (m *Point) Validate() error {
	if m.Longitude < MinimumLongitude || m.Longitude > MaximumLongitude {
		return ValidationError{"invalid longitude value"}
	}

	if m.Latitude < MinimumLatitude || m.Latitude > MaximumLatitude {
		return ValidationError{"invalid latitude value"}
	}

	return nil
}
