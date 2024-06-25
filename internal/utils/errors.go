package utils

import "errors"

var (
	ErrNotExist           = errors.New("requested data does not exist")
	InvalidForm           = errors.New("invalid form")
	ErrEncoding           = errors.New("error encoding message")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrBookAlreadyExists  = errors.New("book already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrBadInput           = errors.New("invalid input")
)
