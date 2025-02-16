package common

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrBadRequest = errors.New("bad request")
	ErrInternal = errors.New("internal error")
	ErrNotFound = errors.New("not found")
)
