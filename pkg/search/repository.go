package search

import "azure.com/ecovo/trip-search-service/pkg/entity"

// Repository is an interface representing the ability to perform CRUD
// operations on searches in a database.
type Repository interface {
	FindByID(ID entity.ID) (*entity.Search, error)
	Create(trip *entity.Search) (entity.ID, error)
	Delete(ID entity.ID) error
}
