package model

import "errors"

// ErrDataProviderAuthentication conveys data pull failed due to authentication
type ErrDataProviderAuthentication struct{}

func (ErrDataProviderAuthentication) Error() string {
	return "unauthenticated"
}

// ErrDataProviderAuthorization conveys data pull failed due to authorization
type ErrDataProviderAuthorization struct{}

func (ErrDataProviderAuthorization) Error() string {
	return "unauthorized"
}

// IsAuthError returns true if this is an authorization or authentication error
func IsAuthError(err error) bool {
	return errors.Is(err, ErrDataProviderAuthorization{}) ||
		errors.Is(err, ErrDataProviderAuthentication{})
}
