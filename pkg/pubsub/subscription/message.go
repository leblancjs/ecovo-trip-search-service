package subscription

// A Message represents data that can be published or received on a
// subscription.
type Message struct {
	// Type represents what the message is about.
	Type string
	// Data represents the message's content.
	Data interface{}
}
