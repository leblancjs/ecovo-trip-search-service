package trip

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"azure.com/ecovo/trip-search-service/pkg/entity"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

// A RestRepository is a repository that performs HTTP requests on trips from the trip-service.
type RestRepository struct {
	domain    string
	authToken string
}

type document struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	DriverID    primitive.ObjectID `bson:"driverId"`
	Vehicle     *entity.Vehicle    `bson:"vehicle"`
	Source      *entity.Point      `bson:"source"`
	Destination *entity.Point      `bson:"destination"`
	LeaveAt     time.Time          `bson:"leaveAt"`
	ArriveBy    time.Time          `bson:"arriveBy"`
	Seats       int                `bson:"seats"`
	Stops       []*entity.Point    `bson:"stops"`
	Details     *entity.Details    `bson:"details"`
}

const (
	// LeaveAtString is a string used for query params
	LeaveAtString = "leaveAt"

	// ArriveByString is a string used for query params
	ArriveByString = "arriveBy"
)

// NewRestRepository creates a trip repository for a MongoDB collection.
func NewRestRepository(domain string, authToken string) (Repository, error) {
	if domain == "" {
		return nil, fmt.Errorf("trip.Rest Repository: domain is nil")
	}

	if authToken == "" {
		return nil, fmt.Errorf("trip.Rest Repository: authToken is nil")
	}

	return &RestRepository{domain, authToken}, nil
}

// Find retrieves all trips based on given filters.
func (r *RestRepository) Find(f *entity.Filters) ([]*entity.Trip, error) {
	req, err := http.NewRequest("GET", "https://"+r.domain+"/trips", nil)
	if err != nil {
		return nil, UnauthorizedError{fmt.Sprintf("trip.repository: failed to create request (%s)", err)}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", r.authToken))

	q := req.URL.Query()

	params, err := f.ToMap()
	for key, value := range params {
		if key != LeaveAtString && key != ArriveByString {
			q.Set(key, value)
		}
	}

	req.URL.RawQuery = q.Encode()

	if len(params[LeaveAtString]) > 0 {
		req.URL.RawQuery = fmt.Sprintf("%s%s%s", req.URL.RawQuery, "&"+LeaveAtString+"=", params[LeaveAtString])
	}

	if len(params[ArriveByString]) > 0 {
		req.URL.RawQuery = fmt.Sprintf("%s%s%s", req.URL.RawQuery, "&"+ArriveByString+"=", params[ArriveByString])
	}

	fmt.Printf(req.URL.RawQuery)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("trip.repository: failed request to trip-api (HTTP %d)", resp.StatusCode)
	}

	var trips []*entity.Trip
	err = json.NewDecoder(resp.Body).Decode(&trips)
	if err != nil {
		return nil, fmt.Errorf("trip.repository: failed to decode trips (%s)", err)
	}
	return trips, nil
}
