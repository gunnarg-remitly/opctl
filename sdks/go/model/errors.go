package model

// ErrDataProviderAuthentication conveys data pull failed due to authentication
type ErrDataProviderAuthentication struct{}

func (ear ErrDataProviderAuthentication) Error() string {
	return "unauthenticated"
}

// ErrDataProviderAuthorization conveys data pull failed due to authorization
type ErrDataProviderAuthorization struct{}

func (ear ErrDataProviderAuthorization) Error() string {
	return "unauthorized"
}
