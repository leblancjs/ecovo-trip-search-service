package trip

import (
	"time"

	"azure.com/ecovo/trip-search-service/pkg/entity"
)

// A MockRepository is a mock repository that returns mocked trips.
type MockRepository struct {
}

// NewMockRepository creates a mock trip repository.
func NewMockRepository() (Repository, error) {
	return &MockRepository{}, nil
}

// Find retrieves all trips based on given filters.
func (r *MockRepository) Find(f *entity.Filters) ([]*entity.Trip, error) {
	var trips []*entity.Trip

	for i := 0; i < 10; i++ {
		trips = append(trips, createMockTrip())
	}

	return trips, nil
}

// createMockTrip creates a mocked trip
func createMockTrip() *entity.Trip {
	return &entity.Trip{
		Source: &entity.Point{
			Latitude:  45.50169,
			Longitude: -73.567253,
			Name:      "Montreal",
		},
		Destination: &entity.Point{
			Latitude:  46.1240621,
			Longitude: -73.5055304,
			Name:      "Quebec",
		},
		LeaveAt:  time.Now().Add(time.Hour * 5),
		ArriveBy: time.Now().Add(time.Hour * 10),
		Seats:    3,
		Stops: []*entity.Point{
			&entity.Point{
				Latitude:  45.9196296,
				Longitude: -73.6149923,
				Name:      "Drummondville",
			},
		},
		Details: &entity.Details{
			Animals:  1,
			Luggages: 1,
		},
	}
}
