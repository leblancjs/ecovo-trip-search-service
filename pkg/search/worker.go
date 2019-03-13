package search

import (
	"fmt"
	"log"
	"time"

	"azure.com/ecovo/trip-search-service/pkg/pubsub/subscription"
)

type Worker struct {
	sub     subscription.Subscription
	started bool
	quit    chan bool
}

func NewWorker(sub subscription.Subscription) (*Worker, error) {
	if sub == nil {
		return nil, fmt.Errorf("search.Worker: cannot work with nil subscription")
	}

	return &Worker{
		sub:  sub,
		quit: make(chan bool),
	}, nil
}

func (w *Worker) Start() {
	if w.started {
		return
	}

	w.started = true

	go w.run()
}

func (w *Worker) Stop() {
	if !w.started {
		return
	}

	w.quit <- true

	w.started = false
}

func (w *Worker) run() {
	i := 0
	for {
		select {
		case <-w.quit:
			return
		default:
			type t struct {
				Topic   string `json:"topic"`
				Message string `json:"message"`
			}

			var err error
			if i%5 == 0 {
				err = w.sub.Publish(&subscription.Event{
					Type: "clearResults",
					Data: nil,
				})
			} else {
				err = w.sub.Publish(&subscription.Event{
					Type: "result",
					Data: t{w.sub.Topic(), fmt.Sprintf("Hello %d", i)},
				})
			}
			if err != nil {
				log.Println(err)
				break
			}

			i++

			time.Sleep(3 * time.Second)
		}
	}
}
