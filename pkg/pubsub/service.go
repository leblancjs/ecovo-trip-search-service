package pubsub

import "azure.com/ecovo/trip-search-service/pkg/pubsub/subscription"

type UseCase interface {
	Subscribe(topic string) (subscription.Subscription, error)
	Unsubscribe(topic string)
}

type Service struct {
	repo subscription.Repository
}

func NewService(repo subscription.Repository) UseCase {
	return &Service{repo}
}

func (s *Service) Subscribe(topic string) (subscription.Subscription, error) {
	return s.repo.Create(topic)
}

func (s *Service) Unsubscribe(topic string) {
	s.repo.Delete(topic)
}
