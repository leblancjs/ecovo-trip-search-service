package entity

import (
	"fmt"
	"strconv"
	"time"
)

// Search contains a search's information.
type Search struct {
	ID      ID       `json:"id"`
	Filters *Filters `json:"filters"`
}

// Validate validates that the search's required fields are filled out
// correctly.
func (s *Search) Validate() error {
	if s.Filters == nil {
		return ValidationError{"missing filters"}
	}

	if err := s.Filters.Validate(); err != nil {
		return err
	}

	return nil
}

// Filters represent the criteria to use to search for trips.
type Filters struct {
	Seats        *int      `json:"seats,ommitempty"`
	LeaveAt      time.Time `json:"leaveAt,ommitempty"`
	ArriveBy     time.Time `json:"arriveBy,ommitempty"`
	Details      *Details  `json:"details,ommitempty"`
	RadiusThresh *int      `json:"radiusThresh,ommitempty"`
	Source       *Point    `json:"source,ommitempty"`
	Destination  *Point    `json:"destination,ommitempty"`
}

const (
	// MinimumRating represents the minimum value a rating can be given
	MinimumRating = 1

	// MaximumRating represents the maximum value a rating can be given
	MaximumRating = 5

	// MinimumLuggagesValue represents the minimum value for luggages
	MinimumLuggagesValue = 0

	// MaximumLuggagesValue represents the maximum value for luggages
	MaximumLuggagesValue = 2

	// MinimumAnimalsValue represents the minimum value for animals
	MinimumAnimalsValue = 0

	// MaximumAnimalsValue represents the maximum value for animals
	MaximumAnimalsValue = 1

	// MinimumRadiusThresh represents the minimum value for radius threshold
	MinimumRadiusThresh = 0

	// FormatTimeRgx is used to convert time to string
	FormatTimeRgx = "2006-01-02T15:04:05Z07:00"
)

// Validate validates that the filters's required fields are filled out correctly.
func (f *Filters) Validate() error {
	if f.Seats != nil && *f.Seats < 0 {
		return ValidationError{"seats filter must be greater than 0"}
	}

	if !f.LeaveAt.IsZero() && !f.ArriveBy.IsZero() {
		return ValidationError{"can't have leaveAt and arriveBy filter at the same time"}
	}

	if f.RadiusThresh != nil && *f.RadiusThresh <= MinimumRadiusThresh {
		return ValidationError{fmt.Sprintf("radiusThresh must be greater than %d", MinimumRadiusThresh)}
	}

	if f.Source != nil {
		err := f.Source.Validate()
		if err != nil {
			return err
		}
	}

	if f.Destination != nil {
		err := f.Destination.Validate()
		if err != nil {
			return err
		}
	}

	if f.Details != nil {
		err := f.Destination.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// ToMap returns list of query params
func (f *Filters) ToMap() (map[string]string, error) {
	mapArr := make(map[string]string)

	if f.Seats != nil {
		mapArr["seats"] = strconv.Itoa(*f.Seats)
	}

	if !f.LeaveAt.IsZero() {
		mapArr["leaveAt"] = f.LeaveAt.Format(FormatTimeRgx)
	}

	if !f.ArriveBy.IsZero() {
		mapArr["arriveBy"] = f.ArriveBy.Format(FormatTimeRgx)
	}

	if f.RadiusThresh != nil {
		mapArr["radiusThresh"] = strconv.Itoa(*f.RadiusThresh)
	}

	if f.Source != nil {
		mapArr["sourceLongitude"] = fmt.Sprintf("%f", f.Source.Longitude)
		mapArr["sourceLatitude"] = fmt.Sprintf("%f", f.Source.Latitude)
	}

	if f.Destination != nil {
		mapArr["destinationLongitude"] = fmt.Sprintf("%f", f.Destination.Longitude)
		mapArr["destinationLatitude"] = fmt.Sprintf("%f", f.Destination.Latitude)
	}

	if f.Details != nil {
		mapArr["detailsAnimals"] = strconv.Itoa(f.Details.Animals)
		mapArr["detailsLuggages"] = strconv.Itoa(f.Details.Luggages)
	}

	return mapArr, nil
}
