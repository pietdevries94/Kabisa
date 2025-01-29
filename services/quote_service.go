package services

import (
	"github.com/pietdevries94/Kabisa/models"
	"github.com/pietdevries94/Kabisa/repositories"
)

type QuoteService interface {
	GetRandomQuote() (*models.Quote, error)
}

type quoteService struct {
	dummyJsonRepo repositories.DummyJsonRepo
}

func NewQuoteService(dummyJsonRepo repositories.DummyJsonRepo) QuoteService {
	return &quoteService{
		dummyJsonRepo: dummyJsonRepo,
	}
}

func (service *quoteService) GetRandomQuote() (*models.Quote, error) {
	return service.dummyJsonRepo.GetRandomQuote()
}
