package services

import "github.com/pietdevries94/Kabisa/models"

type dummyJsonRepo interface {
	GetRandomQuotes(amount int) ([]*models.Quote, error)
}

type quoteGameRepo interface {
	CreateQuoteGame(quotes []*models.Quote) (*models.QuoteGame, error)
}
