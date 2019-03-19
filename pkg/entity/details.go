package entity

import "fmt"

// Details contains a trip's details when it comes to animals,
// and luggages.
type Details struct {
	Animals  int `json:"animals"`
	Luggages int `json:"luggages"`
}

const (
	// AnimalsDetailsNo represents that a driver does not accept animals.
	AnimalsDetailsNo = 0

	// AnimalsDetailsYes represents that a driver accepts animals.
	AnimalsDetailsYes = 1

	// LuggagesDetailsSmall represents that a driver accepts small sized luggages.
	LuggagesDetailsSmall = 0

	// LuggagesDetailsMedium represents that a driver accepts medium sized luggages.
	LuggagesDetailsMedium = 1

	// LuggagesDetailsBig represents that a driver accepts big sized luggages.
	LuggagesDetailsBig = 2
)

// Validate validates that the preferences' required fields are filled out
// correctly.
func (d *Details) Validate() error {
	if d.Animals < AnimalsDetailsNo || d.Animals > AnimalsDetailsYes {
		return ValidationError{fmt.Sprintf("smoking preference is out of bounds \"%d\"", d.Animals)}
	}

	if d.Luggages < LuggagesDetailsSmall || d.Luggages > LuggagesDetailsBig {
		return ValidationError{fmt.Sprintf("conversation preference is out of bounds \"%d\"", d.Luggages)}
	}

	return nil
}
