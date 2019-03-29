package entity

import (
	"time"
)

// Stop struct
type Stop struct {
	ID        ID        `json:"id"`
	Point     *Point    `json:"point,ommitempty"`
	Seats     int       `json:"seats,ommitempty"`
	TimeStamp time.Time `json:"timestamp,ommitempty"`
}

// Validate validates that the stop's required fields are filled out correctly.
func (s *Stop) Validate() error {
	if s.Point != nil {
		err := s.Point.Validate()
		if err != nil {
			return err
		}
	} else {
		return ValidationError{"point is missing"}
	}

	return nil
}
