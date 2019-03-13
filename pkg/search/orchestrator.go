package search

import (
	"fmt"

	"azure.com/ecovo/trip-search-service/pkg/pubsub/subscription"
)

type Orchestrator struct {
	workers map[string]*Worker
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{make(map[string]*Worker)}
}

func (o *Orchestrator) StartWorker(id string, sub subscription.Subscription) error {
	if id == "" {
		return fmt.Errorf("search.Orchestrator: cannot start worker for empty search ID")
	}

	_, ok := o.workers[id]
	if ok {
		return fmt.Errorf("search.Orchestrator: cannot start another worker for same search ID \"%s\"", id)
	}

	worker, err := NewWorker(sub)
	if err != nil {
		return err
	}

	o.workers[id] = worker

	worker.Start()

	return nil
}

func (o *Orchestrator) StopWorker(id string) {
	if worker, ok := o.workers[id]; ok {
		worker.Stop()
	}
}
