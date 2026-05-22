package repository

import (
	"errors"
)

const (
	ErrUniqueViolation = "23505"
	ErrForeignKeeyViolation = "23503"
	ErrNotNullViolation = "23502"
)
var ( 
	ErrAlreadyApplied = errors.New("already applied")
	ErrDuplicate = errors.New("duplicate entry")
	ErrNotFound = errors.New("not found")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)