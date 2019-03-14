package subscription

// A Subscription is an interface representing the ability to publish and
// listen for messages on a given topic.
type Subscription interface {
	Publish(msg *Message) error
	Topic() string
}
