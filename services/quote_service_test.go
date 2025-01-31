package services

import (
	"os"
	"testing"

	"github.com/pietdevries94/Kabisa/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockedDummyJsonRepo struct {
	mock.Mock
}

// Get is fully mocked here, returning with
func (m *MockedDummyJsonRepo) GetRandomQuote() (*models.Quote, error) {
	args := m.Called()
	return args.Get(0).(*models.Quote), args.Error(1)
}

func TestQuoteServiceGetRandomQuote(t *testing.T) {
	t.Run("returns quote from dummyJsonRepo", func(t *testing.T) {
		mockedDummyJsonRepo := new(MockedDummyJsonRepo)

		// The mocked quote is based on an actual response from the underlying api
		mockedDummyJsonRepo.On("GetRandomQuote").Once().Return(&models.Quote{
			ID:     663,
			Quote:  "Never Mistake Motion For Action.",
			Author: "Ernest Hemingway",
		}, nil)

		// We inject the mocked repo into the service and expect the same quote back
		logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
		res, err := NewQuoteService(&logger, mockedDummyJsonRepo).GetRandomQuote()
		require.NoError(t, err)
		assert.Equal(t, &models.Quote{
			ID:     663,
			Quote:  "Never Mistake Motion For Action.",
			Author: "Ernest Hemingway",
		}, res)
	})
	// TODO extend tests and convert to closure tests
}
