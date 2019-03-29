package subscription

import (
	"fmt"
	"strings"

	"github.com/ably/ably-go/ably"
)

// An AblyRepository is a repository that performs CRUD operations on
// subscriptions in an in-memory database. It manages subscriptions to Ably
// realtime channels.
type AblyRepository struct {
	client              *ably.RealtimeClient
	subcriptionsByTopic map[string][]Subscription
}

// NewAblyRepository creates a subscription repository to manage Ably realtime
// channels.
func NewAblyRepository(client *ably.RealtimeClient) (Repository, error) {
	if client == nil {
		return nil, fmt.Errorf("pubsub.AblyRepository: client is nil")
	}

	return &AblyRepository{
		client:              client,
		subcriptionsByTopic: make(map[string][]Subscription),
	}, nil
}

// Create creates a new Ably realtime channel for the given topic, and a
// subscription to manage it.
func (r *AblyRepository) Create(topic string) (Subscription, error) {
	if topic == "" {
		return nil, fmt.Errorf("subscription.AblyRepository: topic cannot be empty")
	}

	channel := r.client.Channels.Get(topic)
	sub, err := NewAblySubscription(channel)
	if err != nil {
		return nil, err
	}

	r.addSubscriptionToTopic(sub)

	return sub, nil
}

// Delete destroys the subscription for the given topic, destroying the Ably
// realtime channel it manages.
func (r *AblyRepository) Delete(topic string) {
	r.removeSubscriptionFromTopic(topic)
}

func (r *AblyRepository) addSubscriptionToTopic(sub Subscription) error {
	if sub == nil {
		return fmt.Errorf("pubsub.Service: cannot add nil subscription to topic")
	}

	topic := sub.Topic()

	subs, ok := r.subcriptionsByTopic[topic]
	if !ok {
		subs = make([]Subscription, 0, 1)
	}

	r.subcriptionsByTopic[topic] = append(subs, sub)

	return nil
}

func (r *AblyRepository) removeSubscriptionFromTopic(topic string) {
	if topic == "" {
		return
	}

	subs, ok := r.subcriptionsByTopic[topic]
	if !ok {
		return
	}

	for i, sub := range subs {
		if strings.Compare(sub.Topic(), topic) == 0 {
			r.subcriptionsByTopic[topic] = append(subs[:i], subs[i+1:]...)
			break
		}
	}
}
