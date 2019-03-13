package subscription

import (
	"encoding/json"
	"fmt"

	"github.com/ably/ably-go/ably"
)

// An AblySubscription represents a subscription to an Ably realtime channel.
type AblySubscription struct {
	// Channel represents the Ably realtime channel associated to the
	// subscription's topic.
	channel *ably.RealtimeChannel
}

// NewAblySubscription creates a new subscription to the given channel.
func NewAblySubscription(channel *ably.RealtimeChannel) (Subscription, error) {
	if channel == nil {
		return nil, fmt.Errorf("subscription.AblySubscription: channel cannot be nil")
	}

	return &AblySubscription{channel}, nil
}

func (s *AblySubscription) Publish(event *Event) error {
	if event == nil {
		return fmt.Errorf("subscription.AblySubscription [topic=%s]: event cannot be nil", s.Topic())
	}

	payload, err := json.Marshal(event.Data)
	if err != nil {
		return fmt.Errorf("subscription.AblySubscription [topic=%s]: failed to marshal message (%s)", s.Topic(), err)
	}

	res, err := s.channel.Publish(event.Type, string(payload))
	if err != nil {
		return fmt.Errorf("subscription.AblySubscription [topic=%s]: failed to publish message (%s)", s.Topic(), err)
	}

	err = res.Wait()
	if err != nil {
		return fmt.Errorf("subscription.AblySubscription [topic=%s]: failed to marshal message (%s)", s.Topic(), err)
	}

	return nil
}

// Topic returns the subscriptions's topic.
func (s *AblySubscription) Topic() string {
	return s.channel.Name
}
