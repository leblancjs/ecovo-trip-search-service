package search

import (
	"fmt"

	"azure.com/ecovo/trip-search-service/pkg/entity"
	"azure.com/ecovo/trip-search-service/pkg/pubsub"
)

// UseCase is an interface representing the ability to handle the business
// logic that involves searches for trips.
type UseCase interface {
	Create(search *entity.Search) (*entity.Search, error)
	FindByID(ID entity.ID) (*entity.Search, error)
	Delete(ID entity.ID) error
}

// A Service handles the business logic related to searches for trips.
type Service struct {
	repo         Repository
	pubSub       pubsub.UseCase
	orchestrator *Orchestrator
}

// NewService creates a search service to handle business logic and manipulate
// searches through a repository.
func NewService(repo Repository, pubSub pubsub.UseCase) UseCase {
	return &Service{repo, pubSub, NewOrchestrator()}
}

// Create validates the search's information, creates it, creates a
// subscription and starts searching for results in the background that will be
// published to the subscription.
func (s *Service) Create(search *entity.Search) (*entity.Search, error) {
	if search == nil {
		return nil, fmt.Errorf("trip.Service: trip is nil")
	}

	err := search.Validate()
	if err != nil {
		return nil, err
	}

	search.ID, err = s.repo.Create(search)
	if err != nil {
		return nil, err
	}

	sub, err := s.pubSub.Subscribe(string(search.ID))
	if err != nil {
		return nil, err
	}

	err = s.orchestrator.StartSearch(search, sub)
	if err != nil {
		s.pubSub.Unsubscribe(search.ID.Hex())
		_ = s.repo.Delete(search.ID)
		return nil, err
	}

	return search, nil
}

// FindByID retrieves the search with the given ID in the repository, if it
// exists.
func (s *Service) FindByID(ID entity.ID) (*entity.Search, error) {
	search, err := s.repo.FindByID(ID)
	if err != nil {
		return nil, NotFoundError{err.Error()}
	}

	return search, nil
}

// Delete erases the search from the repository, and stops searching for
// results.
func (s *Service) Delete(ID entity.ID) error {
	s.orchestrator.StopSearch(ID.Hex())

	s.pubSub.Unsubscribe(string(ID))

	err := s.repo.Delete(ID)
	if err != nil {
		return err
	}

	return nil
}
