package contracts

import "errors"

var (
	// ErrInvalidCredentials is returned when email/password do not match.
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrTokenExpired is returned when the JWT token has expired.
	ErrTokenExpired = errors.New("token expired")

	// ErrTokenInvalid is returned when the JWT token cannot be parsed or is structurally invalid.
	ErrTokenInvalid = errors.New("invalid token")
)
