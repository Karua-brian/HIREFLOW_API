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
	ErrUserEmailExists = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrRecruiterRequestAlreadyExists = errors.New("a pending recruiter request already exists for this user")
	ErrInvalidStatus = errors.New("invalid status value, must be 'approved' or 'rejected'")
	ErrRecruiterRequestNotFound = errors.New("recruiter request not found")
	ErrAlreadyAppliedRequest = errors.New("request already processed")
)