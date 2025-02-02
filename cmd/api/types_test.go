package main

import (
	"context"

	"github.com/pietdevries94/Kabisa/models"
	"github.com/stretchr/testify/mock"
)

type MockedQuoteService struct {
	mock.Mock
}

// GetRandomQuote is fully mocked here
func (m *MockedQuoteService) GetRandomQuote(_ context.Context) (*models.Quote, error) {
	args := m.Called()
	return args.Get(0).(*models.Quote), args.Error(1)
}

// CreateQuoteGame is fully mocked here
func (m *MockedQuoteService) CreateQuoteGame(_ context.Context) (*models.QuoteGame, error) {
	args := m.Called()
	return args.Get(0).(*models.QuoteGame), args.Error(1)
}
