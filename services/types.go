package services

import (
	"context"

	"github.com/pietdevries94/Kabisa/models"
)

type dummyJsonRepo interface {
	GetRandomQuotes(ctx context.Context, amount int) ([]*models.Quote, error)
}

type quoteGameRepo interface {
	CreateQuoteGame(ctx context.Context, quotes []*models.Quote) (*models.QuoteGame, error)
}
