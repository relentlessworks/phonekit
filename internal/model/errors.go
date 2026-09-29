package model

import "errors"

var (
	ErrEmptyNumber   = errors.New("empty phone number")
	ErrTooShort      = errors.New("phone number too short (minimum 4 digits)")
	ErrTooLong       = errors.New("phone number too long (maximum 15 digits)")
	ErrInvalidFormat = errors.New("invalid format specified")
)
