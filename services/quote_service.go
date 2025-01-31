package services

import (
	"github.com/pietdevries94/Kabisa/models"
	"github.com/pietdevries94/Kabisa/repositories"
	"github.com/rs/zerolog"
)

type QuoteService interface {
	GetRandomQuote() (*models.Quote, error)
}

type quoteService struct {
	logger        *zerolog.Logger
	dummyJsonRepo repositories.DummyJsonRepo
}

func NewQuoteService(logger *zerolog.Logger, dummyJsonRepo repositories.DummyJsonRepo) QuoteService {
	return &quoteService{
		logger:        logger,
		dummyJsonRepo: dummyJsonRepo,
	}
}

func (service *quoteService) GetRandomQuote() (*models.Quote, error) {
	return service.dummyJsonRepo.GetRandomQuote()
}
