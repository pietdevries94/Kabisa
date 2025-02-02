package main

import "github.com/pietdevries94/Kabisa/models"

type quoteService interface {
	GetRandomQuote() (*models.Quote, error)
	CreateQuoteGame() (*models.QuoteGame, error)
}
