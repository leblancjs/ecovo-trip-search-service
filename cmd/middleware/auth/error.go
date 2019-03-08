package auth

// An UnauthorizedError is an error that occurs when the user's authorization
// could not be validated.
type UnauthorizedError struct {
	msg string
}

func (e UnauthorizedError) Error() string {
	return e.msg
}
