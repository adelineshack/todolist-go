package core_errors

import "errors"

var (
	ErrNotFound        = errors.New("not_found")
	ErrInvalidArgument = errors.New("invalid argument")
	Conflict           = errors.New("conflict")
)
