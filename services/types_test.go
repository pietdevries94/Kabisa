package services

import (
	"context"

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

type MockedQuoteGameRepo struct {
	mock.Mock
}

func (m *MockedQuoteGameRepo) CreateQuoteGame(_ context.Context, quotes []*models.Quote) (*models.QuoteGame, error) {
	args := m.Called(quotes)
	return args.Get(0).(*models.QuoteGame), args.Error(1)
}
