package errs

import "errors"

var (
	ErrEmailAlreadyExists     = errors.New("email already exists")
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
)
