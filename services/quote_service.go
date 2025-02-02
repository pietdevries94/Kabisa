package services

import (
	"errors"

	"github.com/pietdevries94/Kabisa/models"
	"github.com/rs/zerolog"
)

type QuoteService struct {
	logger        *zerolog.Logger
	dummyJsonRepo dummyJsonRepo
	quoteGameRepo quoteGameRepo
}

func NewQuoteService(logger *zerolog.Logger, dummyJsonRepo dummyJsonRepo, quoteGameRepo quoteGameRepo) *QuoteService {
	return &QuoteService{
		logger:        logger,
		dummyJsonRepo: dummyJsonRepo,
		quoteGameRepo: quoteGameRepo,
	}
}

func (service *QuoteService) GetRandomQuote() (*models.Quote, error) {
	res, err := service.dummyJsonRepo.GetRandomQuotes(1)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("dummyJsonRepo returned no quotes and no error")
	}
	return res[0], nil
}

func (service *QuoteService) CreateQuoteGame() (*models.QuoteGame, error) {
	quotes, err := service.dummyJsonRepo.GetRandomQuotes(3)
	if err != nil {
		return nil, err
	}
	return service.quoteGameRepo.CreateQuoteGame(quotes)
}
