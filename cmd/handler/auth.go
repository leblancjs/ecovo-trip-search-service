package handler

import (
	"context"
	"net/http"

	"azure.com/ecovo/trip-search-service/cmd/middleware/auth"
)

// Auth validates a request's authorization header using the given validator
// to ensure that the user is authorized to access an endpoint and extracts the
// authenticated user's information.
//
// The authenticated user's information placed in the request's context and can
// be accessed by using the auth.FromContext utility function.
func Auth(validator auth.Validator, next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		userInfo, err := validator.Validate(r.Header.Get("Authorization"))
		if err != nil {
			return err
		}

		ctx := context.WithValue(r.Context(), auth.UserInfoContextKey, userInfo)
		next.ServeHTTP(w, r.WithContext(ctx))

		return nil
	}
}
