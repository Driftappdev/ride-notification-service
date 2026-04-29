package error

import "errors"

var (
	ErrInvalidLocation = errors.New("invalid location data")
	ErrDriverNotFound  = errors.New("driver not found")
	ErrInternal        = errors.New("internal server error")
)
