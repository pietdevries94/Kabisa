package main

import (
	"github.com/pietdevries94/Kabisa/models"
	"github.com/stretchr/testify/mock"
)

type MockedQuoteService struct {
	mock.Mock
}

// GetRandomQuote is fully mocked here
func (m *MockedQuoteService) GetRandomQuote() (*models.Quote, error) {
	args := m.Called()
	return args.Get(0).(*models.Quote), args.Error(1)
}

// CreateQuoteGame is fully mocked here
func (m *MockedQuoteService) CreateQuoteGame() (*models.QuoteGame, error) {
	args := m.Called()
	return args.Get(0).(*models.QuoteGame), args.Error(1)
}
