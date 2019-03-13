package subscription

type Repository interface {
	Create(topic string) (Subscription, error)
	Delete(topic string)
}
