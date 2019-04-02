package search

import (
	"fmt"

	"azure.com/ecovo/trip-search-service/pkg/entity"
	"azure.com/ecovo/trip-search-service/pkg/pubsub/subscription"
	"azure.com/ecovo/trip-search-service/pkg/route"
)

// An Orchestrator manages workers that run asynchronously to gather search
// results and publish them on subscriptions. It creates, starts, stops and
// deletes them.
type Orchestrator struct {
	workers      map[string]*Worker
	routeService route.UseCase
}

// NewOrchestrator creates a search orchestrator to manage workers that run to
// search for results asynchronously.
func NewOrchestrator(routeService route.UseCase) *Orchestrator {
	return &Orchestrator{make(map[string]*Worker), routeService}
}

// StartSearch creates and starts a worker to search for results and publish
// them to the given subscription. Only one worker can exist for a given search
// ID.
func (o *Orchestrator) StartSearch(search *entity.Search, sub subscription.Subscription) error {
	if search == nil {
		return fmt.Errorf("search.Orchestrator: cannot start worker for nil search")
	}

	searchID := search.ID.Hex()

	_, ok := o.workers[searchID]
	if ok {
		return fmt.Errorf("search.Orchestrator: cannot start another worker for same search ID \"%s\"", searchID)
	}

	worker, err := NewWorker(search.Filters, sub, o.routeService)
	if err != nil {
		return err
	}

	o.workers[searchID] = worker

	worker.Start()

	return nil
}

// StopSearch stops a worker, stopping the search.
func (o *Orchestrator) StopSearch(id string) {
	if id == "" {
		return
	}

	if worker, ok := o.workers[id]; ok {
		worker.Stop()
	}

	delete(o.workers, id)
}

// PublishTrip will send trips to the work so it can be published on Ably
func (o *Orchestrator) PublishTrip(trip *entity.Trip) {
	for _, w := range o.workers {
		w.trips <- trip
	}
}
