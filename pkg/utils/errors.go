package utils

import "errors"

var (
	ErrorQueryString       = errors.New("Error create query string")
	ErrorUserAlreadyExists = errors.New("Username or email already exists")
	ErrorUserNotFound      = errors.New("User not found")
	ErrorInvalidPassword   = errors.New("Invalid password")
)
