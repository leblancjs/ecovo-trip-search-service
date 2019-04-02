package subscription

// Callback is a callback function for the subcsription
type Callback func(msg *Message)

// A Subscription is an interface representing the ability to publish and
// listen for messages on a given topic.
type Subscription interface {
	Publish(msg *Message) error
	Subscribe(callback Callback) error
	Topic() string
}
