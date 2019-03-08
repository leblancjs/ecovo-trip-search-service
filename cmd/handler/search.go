package handler

import (
	"encoding/json"
	"net/http"

	"azure.com/ecovo/trip-search-service/pkg/entity"
	"azure.com/ecovo/trip-search-service/pkg/search"
	"github.com/gorilla/mux"
)

// StartSearch handles a request to start searching for a trip.
func StartSearch(service search.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		var s *entity.Search
		err := json.NewDecoder(r.Body).Decode(&s)
		if err != nil {
			return err
		}

		s, err = service.Create(s)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(s)
		if err != nil {
			_ = service.Delete(entity.ID(s.ID))

			return err
		}

		return nil
	}
}

// GetSearchByID handles a request to retrieve a search by its unique identifier.
func GetSearchByID(service search.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)

		id := entity.NewIDFromHex(vars["id"])
		t, err := service.FindByID(id)
		if err != nil {
			return err
		}

		err = json.NewEncoder(w).Encode(t)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusOK)

		return nil
	}
}

// StopSearch handles a request to stop searching for a trip by its unique
// identifier.
func StopSearch(service search.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)

		id := entity.NewIDFromHex(vars["id"])

		err := service.Delete(id)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusOK)

		return nil
	}
}
