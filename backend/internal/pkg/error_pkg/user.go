package errorpkg

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)
