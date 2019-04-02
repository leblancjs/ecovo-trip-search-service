package search

import (
	"fmt"
	"log"

	"azure.com/ecovo/trip-search-service/pkg/entity"
	"azure.com/ecovo/trip-search-service/pkg/pubsub/subscription"
	"azure.com/ecovo/trip-search-service/pkg/route"
	"github.com/umahmood/haversine"
	"googlemaps.github.io/maps"
)

// A Worker does all the heavy lifting to search for trips that either match
// the search filters or come close. It runs in a Go routine, to avoid blocking
// the entire service, and publishes results to a subscription.
type Worker struct {
	filters      *entity.Filters
	sub          subscription.Subscription
	started      bool
	trips        chan *entity.Trip
	routeService route.UseCase
	quit         chan bool
}

// NewWorker creates a new search worker that uses the subscription to publish
// results.
func NewWorker(filters *entity.Filters, sub subscription.Subscription, routeService route.UseCase) (*Worker, error) {
	if filters == nil {
		return nil, fmt.Errorf("search.Worker: cannot work with nil filters")
	}

	if sub == nil {
		return nil, fmt.Errorf("search.Worker: cannot work with nil subscription")
	}

	return &Worker{
		filters:      filters,
		sub:          sub,
		trips:        make(chan *entity.Trip),
		routeService: routeService,
		quit:         make(chan bool),
	}, nil
}

// Start tells the worker to start searching for trips.
func (w *Worker) Start() {
	if w.started {
		return
	}

	w.started = true

	go w.run()
}

// Stop tells the worker to stop searching for trips.
func (w *Worker) Stop() {
	if !w.started {
		return
	}

	w.quit <- true

	w.started = false
}

func (w *Worker) run() {
	for {
		select {
		case <-w.quit:
			return
		case trip := <-w.trips:
			r, err := w.routeService.GetRoute(trip)
			if err != nil {
				log.Println("searchWorker: failed to get route from google maps")
				break
			}

			isValid := false
			if w.filters != nil && trip != nil {
				isValid, err = validateTrip(trip, w.filters, r)
				if err != nil {
					log.Println(err)
					break
				}
			}

			if isValid {
				err := w.sub.Publish(&subscription.Message{
					Type: EventAddResult,
					Data: trip,
				})
				if err != nil {
					log.Println(err)
					break
				}
			}
		default:
			break
		}
	}
}

// validateTrip will validate
func validateTrip(t *entity.Trip, f *entity.Filters, route maps.Route) (bool, error) {
	threshold := metersToKM(float64(*f.RadiusThresh))
	points, err := route.OverviewPolyline.Decode()
	if err != nil {
		return false, err
	}

	source := haversine.Coord{Lat: f.Source.Latitude, Lon: f.Source.Longitude}
	destination := haversine.Coord{Lat: f.Destination.Latitude, Lon: f.Destination.Longitude}

	isSourceOk := false
	isDestinationOk := false

	for _, p := range points {
		routePoint := haversine.Coord{Lat: p.Lat, Lon: p.Lng}
		_, sourceDistance := haversine.Distance(source, routePoint)
		_, destinationDistance := haversine.Distance(destination, routePoint)

		if sourceDistance <= threshold {
			isSourceOk = true
		}

		if destinationDistance <= threshold {
			isDestinationOk = true
		}
	}

	if isSourceOk && isDestinationOk {
		return true, nil
	}

	return false, nil
}

// metersToKM converts meters into kilometers
func metersToKM(meters float64) float64 {
	return meters / 1000.0
}
