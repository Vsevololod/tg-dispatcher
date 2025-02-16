package storage

import "errors"

var (
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrVideoNotFound = errors.New("video not found")
	ErrAppNotFound   = errors.New("app not found")
)
