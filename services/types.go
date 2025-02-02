package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/pietdevries94/Kabisa/models"
)

type dummyJsonRepo interface {
	GetRandomQuotes(ctx context.Context, amount int) ([]*models.Quote, error)
	GetQuotes(ctx context.Context, ids []int) (map[int]*models.Quote, error)
}

type quoteGameRepo interface {
	CreateQuoteGame(ctx context.Context, quotes []*models.Quote) (*models.QuoteGame, error)
	ValidateIDAndAnswerIDs(ctx context.Context, id uuid.UUID, answers models.QuoteGameAnswerMap) (answerIDs []int, err error)
	ValidateAnswersAndCreateGameResult(ctx context.Context, id uuid.UUID, quotes map[int]*models.Quote, answers models.QuoteGameAnswerMap) (*models.QuoteGameResult, error)
}
