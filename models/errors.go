package models

import "fmt"

// PublicError is the only type of error that gets returned by the api. No other type of error should be returned.
type PublicError struct {
	msg string
}

func (pe *PublicError) Error() string {
	return pe.msg
}

var _ error = &PublicError{}

// NewPublicError creates a new public error
func NewPublicError(msg string) *PublicError {
	return &PublicError{
		msg: msg,
	}
}

// NewPublicErrorf creates a new public error using fmt.Sprintf to format the message
func NewPublicErrorf(msg string, args ...any) *PublicError {
	return NewPublicError(fmt.Sprintf(msg, args...))
}

var (
	ErrQuoteGameIdNotFound = NewPublicError("quote_game_id_not_found")
	ErrInvalidQuoteID      = NewPublicError("invalid_quote_id")
)
