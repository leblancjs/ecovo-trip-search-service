package handler

import (
	"fmt"
	"net/http"

	"azure.com/ecovo/trip-search-service/cmd/middleware/auth"
)

// An Error is an application error that can be handled by a handler.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   error  `json:"-"`
}

func (err Error) String() string {
	return fmt.Sprintf("code=%d, message=\"%s\", error=\"%s\"", err.Code, err.Message, err.Error)
}

// WrapError wraps the given error in an application error that can be handled
// by a handler.
func WrapError(err error) *Error {
	if err == nil {
		return nil
	} else if _, ok := err.(auth.UnauthorizedError); ok {
		return &Error{http.StatusUnauthorized, "unauthorized", err}
	} else {
		return &Error{
			http.StatusInternalServerError,
			"Something went wrong while processing your request. Please contact your system administrator.",
			err,
		}
	}
}
