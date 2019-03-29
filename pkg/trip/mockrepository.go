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
		LeaveAt:  time.Now().Add(time.Hour * 5),
		ArriveBy: time.Now().Add(time.Hour * 10),
		Seats:    3,
		Stops: []*entity.Stop{
			&entity.Stop{
				ID: entity.NilID,
				Point: &entity.Point{
					Latitude:  45.4944494,
					Longitude: -73.561703,
					Name:      "Montreal",
				},
				Seats:     3,
				TimeStamp: time.Now(),
			},
			&entity.Stop{
				ID: entity.NilID,
				Point: &entity.Point{
					Latitude:  45.881168,
					Longitude: -72.484734,
					Name:      "Drummondville",
				},
				Seats:     3,
				TimeStamp: time.Now(),
			},
			&entity.Stop{
				ID: entity.NilID,
				Point: &entity.Point{
					Latitude:  46.813877,
					Longitude: -71.207977,
					Name:      "Quebec",
				},
				Seats:     3,
				TimeStamp: time.Now(),
			},
		},
		Details: &entity.Details{
			Animals:  1,
			Luggages: 1,
		},
	}
}
