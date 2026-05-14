package domain

import "errors"

var (
	ErrNotFound      = errors.New("record not found")
	ErrConflict      = errors.New("record already exists")
	ErrEntityInUse   = errors.New("cannot delete record because it is currently in use")
	ErrInvalidInput  = errors.New("invalid input data")
	ErrUnauthorized  = errors.New("unauthorized access")
	ErrInternalError = errors.New("internal server error")
)
