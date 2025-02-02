package models

type PublicError struct {
	msg string
}

func (pe *PublicError) Error() string {
	return pe.msg
}

var _ error = &PublicError{}

func NewPublicError(msg string) *PublicError {
	return &PublicError{
		msg: msg,
	}
}

var (
	ErrQuoteGameIdNotFound = NewPublicError("quote_game_id_not_found")
	ErrInvalidQuoteID      = NewPublicError("invalid_quote_id")
)
