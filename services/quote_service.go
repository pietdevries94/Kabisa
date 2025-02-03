package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
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

func (service *QuoteService) GetRandomQuote(ctx context.Context) (*models.Quote, error) {
	res, err := service.dummyJsonRepo.GetRandomQuotes(ctx, 1)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("dummyJsonRepo returned no quotes and no error")
	}
	return res[0], nil
}

func (service *QuoteService) CreateQuoteGame(ctx context.Context) (*models.QuoteGame, error) {
	quotes, err := service.dummyJsonRepo.GetRandomQuotes(ctx, 3)
	if err != nil {
		return nil, err
	}
	return service.quoteGameRepo.CreateQuoteGame(ctx, quotes)
}

func (service *QuoteService) SubmitAnswerToQuoteGame(ctx context.Context, id uuid.UUID, answers models.QuoteGameAnswerMap) (*models.QuoteGameResult, error) {
	quoteIDs, err := service.quoteGameRepo.ValidateIDAndAnswerIDs(ctx, id, answers)
	if err != nil {
		return nil, err
	}

	quotes, err := service.dummyJsonRepo.GetQuotes(ctx, quoteIDs)
	if err != nil {
		return nil, err
	}

	return service.quoteGameRepo.ValidateAnswersAndCreateGameResult(ctx, id, quoteIDs, quotes, answers)
}
