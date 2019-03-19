package trip

// A UnauthorizedError is an error that represents that no trip was found.
type UnauthorizedError struct {
	msg string
}

func (e UnauthorizedError) Error() string {
	return e.msg
}
