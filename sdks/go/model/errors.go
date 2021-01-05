package model

// ErrDataProviderAuthentication conveys data pull failed due to authentication
type ErrDataProviderAuthentication struct{}

func (ear ErrDataProviderAuthentication) Error() string {
	return "Data pull failed due to invalid/lack of authentication"
}

// ErrDataProviderAuthorization conveys data pull failed due to authorization
type ErrDataProviderAuthorization struct{}

func (ear ErrDataProviderAuthorization) Error() string {
	return "Data pull failed due to insufficient/lack of authorization"
}
