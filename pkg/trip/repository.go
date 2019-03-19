package trip

import (
	"azure.com/ecovo/trip-search-service/pkg/entity"
)

// Repository is an interface representing the ability to perform CRUD
// operations on trips in a database.
type Repository interface {
	Find(filters *entity.Filters) ([]*entity.Trip, error)
}
