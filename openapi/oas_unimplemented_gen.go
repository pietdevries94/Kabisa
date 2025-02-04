// Code generated by ogen, DO NOT EDIT.

package openapi

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// CreateNewQuoteGame implements createNewQuoteGame operation.
//
// The quote game returns three quotes and three authors. In `PUT /quote-game/:id`, the player can
// respond with their answer. There is a deadline of five minutes.
//
// POST /quote-game
func (UnimplementedHandler) CreateNewQuoteGame(ctx context.Context) (r CreateNewQuoteGameRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetRandomQuote implements getRandomQuote operation.
//
// Returns a random quote.
//
// GET /quote
func (UnimplementedHandler) GetRandomQuote(ctx context.Context) (r GetRandomQuoteRes, _ error) {
	return r, ht.ErrNotImplemented
}

// SubmitAnswerForQuoteGame implements submitAnswerForQuoteGame operation.
//
// This request expects an answer from the user and will return if the answer was correct and what
// the correct answer should be.
//
// POST /quote-game/{id}/answer
func (UnimplementedHandler) SubmitAnswerForQuoteGame(ctx context.Context, req []QuoteGameAnswer, params SubmitAnswerForQuoteGameParams) (r SubmitAnswerForQuoteGameRes, _ error) {
	return r, ht.ErrNotImplemented
}
