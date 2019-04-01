package entity

import (
	"fmt"
	"time"
)

// Vehicle contains a vehicle's information.
type Vehicle struct {
	ID    ID     `json:"id"`
	Make  string `json:"make"`
	Year  int    `json:"year"`
	Model string `json:"model"`
}

const (
	// YearMinimum represents the minimum year of a car.
	YearMinimum = 1900
)

// Validate validates that the vehicles's required fields are filled out correctly.
func (v *Vehicle) Validate() error {
	if v.ID.IsZero() {
		return ValidationError{fmt.Sprintf("id must not be nil")}
	}

	if v.Year <= YearMinimum && v.Year > time.Now().Year() {
		return ValidationError{fmt.Sprintf("year must me between %d and %d", YearMinimum, time.Now().Year())}
	}

	if v.Make == "" {
		return ValidationError{"make is missing"}
	}

	if v.Model == "" {
		return ValidationError{"model is missing"}
	}

	return nil
}
