package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/pietdevries94/Kabisa/models"
	"github.com/stretchr/testify/mock"
)

type MockedDummyJsonRepo struct {
	mock.Mock
}

func (m *MockedDummyJsonRepo) GetRandomQuotes(_ context.Context, amount int) ([]*models.Quote, error) {
	args := m.Called(amount)
	return args.Get(0).([]*models.Quote), args.Error(1)
}

func (m *MockedDummyJsonRepo) GetQuotes(_ context.Context, ids []int) (map[int]*models.Quote, error) {
	args := m.Called(ids)
	return args.Get(0).(map[int]*models.Quote), args.Error(1)
}

type MockedQuoteGameRepo struct {
	mock.Mock
}

func (m *MockedQuoteGameRepo) CreateQuoteGame(_ context.Context, quotes []*models.Quote) (*models.QuoteGame, error) {
	args := m.Called(quotes)
	return args.Get(0).(*models.QuoteGame), args.Error(1)
}

func (m *MockedQuoteGameRepo) ValidateIDAndAnswerIDs(_ context.Context, id uuid.UUID, answers models.QuoteGameAnswerMap) (quoteIDs []int, err error) {
	args := m.Called(id, answers)
	return args.Get(0).([]int), args.Error(1)
}

func (m *MockedQuoteGameRepo) ValidateAnswersAndCreateGameResult(_ context.Context, id uuid.UUID, quoteIDs []int, quotes map[int]*models.Quote, answers models.QuoteGameAnswerMap) (*models.QuoteGameResult, error) {
	args := m.Called(id, quoteIDs, quotes, answers)
	return args.Get(0).(*models.QuoteGameResult), args.Error(1)
}
