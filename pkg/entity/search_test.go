package entity

import (
	"testing"
)

func TestSearchValidation(t *testing.T) {
	var filters = Filters{}
	var search = Search{Filters: &filters}

	t.Run("Should fail when filters are missing", func(t *testing.T) {
		s := search
		s.Filters = nil

		if _, ok := s.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should succeed when filters are present and valid", func(t *testing.T) {
		s := search

		err := s.Validate()
		if err != nil {
			t.Error(err)
		}
	})
}

func TestFiltersValidation(t *testing.T) {
	var filters = Filters{}

	t.Run("Should succeed when filters are valid", func(t *testing.T) {
		f := filters

		err := f.Validate()
		if err != nil {
			t.Error(err)
		}
	})
}
