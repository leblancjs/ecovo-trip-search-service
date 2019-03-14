package subscription

// A Repository is an interface representing the ability to perform CRUD
// operations on subscriptions.
type Repository interface {
	Create(topic string) (Subscription, error)
	Delete(topic string)
}
