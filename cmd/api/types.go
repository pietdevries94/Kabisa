package main

import (
	"context"

	"github.com/pietdevries94/Kabisa/models"
)

type quoteService interface {
	GetRandomQuote(ctx context.Context) (*models.Quote, error)
	CreateQuoteGame(ctx context.Context) (*models.QuoteGame, error)
}
