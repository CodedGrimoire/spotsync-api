package service

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidRole        = errors.New("invalid role")
	ErrUserNotFound       = errors.New("user not found")
)
