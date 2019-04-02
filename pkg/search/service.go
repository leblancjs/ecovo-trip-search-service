package search

import (
	"encoding/json"
	"fmt"
	"log"

	"azure.com/ecovo/trip-search-service/pkg/entity"
	"azure.com/ecovo/trip-search-service/pkg/pubsub"
	"azure.com/ecovo/trip-search-service/pkg/pubsub/subscription"
	"azure.com/ecovo/trip-search-service/pkg/route"
	"azure.com/ecovo/trip-search-service/pkg/trip"
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
	trip         trip.UseCase
	orchestrator *Orchestrator
}

const searchChannelPrefix = "search:"
const tripsChannel = "trips"

// NewService creates a search service to handle business logic and manipulate
// searches through a repository.
func NewService(repo Repository, pubSub pubsub.UseCase, trip trip.UseCase, routeService route.UseCase) (UseCase, error) {
	orchestrator := NewOrchestrator(routeService)
	tripsSub, err := pubSub.Subscribe(tripsChannel)

	if err != nil {
		return nil, fmt.Errorf("trip.Service: error subscribing to channel (%s) ", err)
	}

	s := &Service{repo, pubSub, trip, orchestrator}

	err = tripsSub.Subscribe(s.listenTripsChange)
	if err != nil {
		return nil, fmt.Errorf("trip.Service: error subscribing to channel (%s) ", err)
	}

	return s, nil
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

	sub, err := s.pubSub.Subscribe(searchChannelPrefix + string(search.ID))
	if err != nil {
		return nil, err
	}

	err = s.orchestrator.StartSearch(search, sub)
	if err != nil {
		s.pubSub.Unsubscribe(searchChannelPrefix + search.ID.Hex())
		_ = s.repo.Delete(search.ID)
		return nil, err
	}

	trips, err := s.trip.Find(search.Filters)
	if err != nil {
		return nil, err
	}

	for _, t := range trips {
		s.orchestrator.PublishTrip(t)
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

	s.pubSub.Unsubscribe(searchChannelPrefix + string(ID))

	err := s.repo.Delete(ID)
	if err != nil {
		return err
	}

	return nil
}

// listenTripsChange is a routine that listens to any update or add of a trip from Ably
func (s *Service) listenTripsChange(msg *subscription.Message) {
	trip := &entity.Trip{}
	err := json.Unmarshal([]byte(msg.Data.(string)), trip)
	if err != nil {
		log.Println("search.Service: unable to unmarshal msg from subscription")
		return
	}
	s.orchestrator.PublishTrip(trip)

}
