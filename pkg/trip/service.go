package trip

import (
	"azure.com/ecovo/trip-search-service/pkg/entity"
)

// UseCase is an interface representing the ability to handle the business
// logic that involves trips.
type UseCase interface {
	Find(filters *entity.Filters) ([]*entity.Trip, error)
}

// A Service handles the business logic related to trips.
type Service struct {
	repo Repository
}

// NewService creates a trip service to handle business logic and manipulate
// trips through a repository.
func NewService(repo Repository) *Service {
	return &Service{repo}
}

// Find retrieves all the trips
func (s *Service) Find(filters *entity.Filters) ([]*entity.Trip, error) {
	err := filters.Validate()
	if err != nil {
		return nil, err
	}

	t, err := s.repo.Find(filters)
	if err != nil {
		return []*entity.Trip{}, err
	}

	return t, nil
}
