package schemas

import "errors"

var (
	ErrRequired      = errors.New("is required")
	ErrMinLength     = errors.New("too short")
	ErrMaxLength     = errors.New("too long")
	ErrInvalidChars  = errors.New("contains invalid characters")
	ErrInvalidAmount = errors.New("invalid amount")
)
