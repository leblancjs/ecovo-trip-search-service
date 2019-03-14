package pubsub

import "azure.com/ecovo/trip-search-service/pkg/pubsub/subscription"

// UseCase is an interface representing the ability to handle the business
// logic that involves subscribing and unsubscribing to topics.
type UseCase interface {
	Subscribe(topic string) (subscription.Subscription, error)
	Unsubscribe(topic string)
}

// A Service handles the business logic related to subscriptions.
type Service struct {
	repo subscription.Repository
}

// NewService creates a pubsub service to handler business logic related to
// subscriptions.
func NewService(repo subscription.Repository) UseCase {
	return &Service{repo}
}

// Subscribe creates a subscription to the given topic.
func (s *Service) Subscribe(topic string) (subscription.Subscription, error) {
	return s.repo.Create(topic)
}

// Unsubscribe destroys a subscription to the given topic.
func (s *Service) Unsubscribe(topic string) {
	s.repo.Delete(topic)
}
