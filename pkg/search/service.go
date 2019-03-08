package search

import (
	"fmt"

	"azure.com/ecovo/trip-search-service/pkg/entity"
)

// UseCase is an interface representing the ability to handle the business
// logic that involves searches for trips.
type UseCase interface {
	Create(t *entity.Search) (*entity.Search, error)
	FindByID(ID entity.ID) (*entity.Search, error)
	Delete(ID entity.ID) error
}

// A Service handles the business logic related to searches for trips.
type Service struct {
	repo Repository
}

// NewService creates a search service to handle business logic and manipulate
// searches through a repository.
func NewService(repo Repository) *Service {
	return &Service{repo}
}

// Create validates the trips's information.
func (s *Service) Create(t *entity.Search) (*entity.Search, error) {
	if t == nil {
		return nil, fmt.Errorf("trip.Service: trip is nil")
	}

	err := t.Validate()
	if err != nil {
		return nil, err
	}

	t.ID, err = s.repo.Create(t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// FindByID retrieves the search with the given ID in the repository, if it
// exists.
func (s *Service) FindByID(ID entity.ID) (*entity.Search, error) {
	t, err := s.repo.FindByID(ID)
	if err != nil {
		return nil, NotFoundError{err.Error()}
	}

	return t, nil
}

// Delete erases the search from the repository.
func (s *Service) Delete(ID entity.ID) error {
	err := s.repo.Delete(ID)
	if err != nil {
		return err
	}

	return nil
}
