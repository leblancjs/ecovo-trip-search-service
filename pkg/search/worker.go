package search

import (
	"fmt"
	"log"

	"azure.com/ecovo/trip-search-service/pkg/entity"
	"azure.com/ecovo/trip-search-service/pkg/pubsub/subscription"
)

// A Worker does all the heavy lifting to search for trips that either match
// the search filters or come close. It runs in a Go routine, to avoid blocking
// the entire service, and publishes results to a subscription.
type Worker struct {
	filters *entity.Filters
	sub     subscription.Subscription
	started bool
	trips   []*entity.Trip
	quit    chan bool
}

// NewWorker creates a new search worker that uses the subscription to publish
// results.
func NewWorker(filters *entity.Filters, sub subscription.Subscription, trips []*entity.Trip) (*Worker, error) {
	if filters == nil {
		return nil, fmt.Errorf("search.Worker: cannot work with nil filters")
	}

	if sub == nil {
		return nil, fmt.Errorf("search.Worker: cannot work with nil subscription")
	}

	return &Worker{
		sub:   sub,
		trips: trips,
		quit:  make(chan bool),
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
		default:
			for t := range w.trips {
				err := w.sub.Publish(&subscription.Message{
					Type: EventAddResult,
					Data: t,
				})
				if err != nil {
					log.Println(err)
					break
				}
			}
			break
		}
	}
}
