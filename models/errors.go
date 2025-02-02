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
