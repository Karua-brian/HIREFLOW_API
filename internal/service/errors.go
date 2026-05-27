package service

import (
	"errors"
)

var (
	ErrBadRequest = errors.New("invalid user id")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden = errors.New("forbidden")
	ErrMessagesUnavailable = errors.New("messages unavailable")
	ErrAlreadyApplied = errors.New("already applied to this job")
	ErrInvalidRole = errors.New("invalid role")
	ErrDuplicate = errors.New("duplicate entry")
	ErrUserExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	
)