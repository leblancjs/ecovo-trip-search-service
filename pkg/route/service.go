package route

import (
	"azure.com/ecovo/trip-search-service/pkg/entity"
	"googlemaps.github.io/maps"
)

// UseCase interface
type UseCase interface {
	GetRoute(t *entity.Trip) (maps.Route, error)
}

// Service structure
type Service struct {
	repo Repository
}

// NewService creates the service
func NewService(repo Repository) UseCase {
	return &Service{repo}
}

// GetRoute returns google maps route for a trip
func (s *Service) GetRoute(t *entity.Trip) (maps.Route, error) {
	return s.repo.GetRoute(t)
}
