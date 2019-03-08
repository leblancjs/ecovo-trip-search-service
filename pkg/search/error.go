package search

// A NotFoundError is an error that represents that no search was found.
type NotFoundError struct {
	msg string
}

func (e NotFoundError) Error() string {
	return e.msg
}
