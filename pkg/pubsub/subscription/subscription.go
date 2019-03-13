package subscription

type Subscription interface {
	Publish(event *Event) error
	Topic() string
}
