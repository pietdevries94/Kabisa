package services

import (
	"github.com/pietdevries94/Kabisa/models"
	"github.com/stretchr/testify/mock"
)

type MockedDummyJsonRepo struct {
	mock.Mock
}

func (m *MockedDummyJsonRepo) GetRandomQuotes(amount int) ([]*models.Quote, error) {
	args := m.Called(amount)
	return args.Get(0).([]*models.Quote), args.Error(1)
}

type MockedQuoteGameRepo struct {
	mock.Mock
}

func (m *MockedQuoteGameRepo) CreateQuoteGame(quotes []*models.Quote) (*models.QuoteGame, error) {
	args := m.Called(quotes)
	return args.Get(0).(*models.QuoteGame), args.Error(1)
}
