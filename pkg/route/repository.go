package route

import (
	"azure.com/ecovo/trip-search-service/pkg/entity"
	"googlemaps.github.io/maps"
)

// Repository interface
type Repository interface {
	GetRoute(t *entity.Trip) (maps.Route, error)
}
