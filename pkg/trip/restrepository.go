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
	VehicleID   primitive.ObjectID `bson:"vehicleId"`
	Source      *entity.Point      `bson:"source"`
	Destination *entity.Point      `bson:"destination"`
	LeaveAt     time.Time          `bson:"leaveAt"`
	ArriveBy    time.Time          `bson:"arriveBy"`
	Seats       int                `bson:"seats"`
	Stops       []*entity.Point    `bson:"stops"`
	Details     *entity.Details    `bson:"details"`
	Steps       []*entity.Point    `bson:"steps"`
}

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
		q.Set(key, value)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("trip.repository: failed to validate token")
	}

	var trips []*entity.Trip
	err = json.NewDecoder(resp.Body).Decode(&trips)
	if err != nil {
		return nil, fmt.Errorf("trip.repository: failed to decode trips (%s)", err)
	}
	return trips, nil
}
