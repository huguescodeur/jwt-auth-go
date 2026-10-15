package user

import "errors"

var (
	ErrAlreadyExist      = errors.New("user already exist")
	ErrNotFound          = errors.New("user not found")
	ErrInvalidCredential = errors.New("invalid credential")
	ErrAllFieldsRequired = errors.New("all fields are requires")
)
