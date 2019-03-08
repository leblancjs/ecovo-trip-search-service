package entity

// Search contains a search's information.
type Search struct {
	ID ID `json:"id"`
}

// Validate validates that the search's required fields are filled out
// correctly.
func (t *Search) Validate() error {
	return nil
}
