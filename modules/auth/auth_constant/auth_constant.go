package authconstant

import "errors"

const (
	ErrCodeInvalidRequest      = "AUTH_INVALID_REQUEST"
	ErrCodeInvalidCredentials  = "AUTH_INVALID_CREDENTIALS"
	ErrCodeLoginFailed         = "AUTH_LOGIN_FAILED"
	ErrCodeDefaultUserSeedFail = "AUTH_DEFAULT_USER_SEED_FAILED"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
