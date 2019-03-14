package entity

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
}

// Validate validates that the filter's required fields are filled out
// correctly.
func (f *Filters) Validate() error {
	return nil
}
