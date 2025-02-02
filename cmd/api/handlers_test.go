package main

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/pietdevries94/Kabisa/api"
	"github.com/pietdevries94/Kabisa/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockedQuoteService struct {
	mock.Mock
}

// Get is fully mocked here, returning with
func (m *MockedQuoteService) GetRandomQuote() (*models.Quote, error) {
	args := m.Called()
	return args.Get(0).(*models.Quote), args.Error(1)
}

func TestGetRandomQuote(t *testing.T) {
	type Test struct {
		mockedServiceQuote *models.Quote
		mockedServiceError error
		expectedResult     api.GetRandomQuoteRes
	}

	run := func(tt Test) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()

			// first we bootstrap a minimal version of the application, needed for the handler
			mockedQuoteService := new(MockedQuoteService)
			mockedQuoteService.On("GetRandomQuote").Once().Return(tt.mockedServiceQuote, tt.mockedServiceError)

			logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
			app := application{
				logger:       &logger,
				quoteService: mockedQuoteService,
			}

			// We now run the handler and validate the result
			res, err := app.GetRandomQuote(context.TODO())
			require.NoError(t, err)
			assert.Equal(t, tt.expectedResult, res)
		}
	}

	t.Run("returns a quote", run(Test{
		mockedServiceQuote: &models.Quote{
			ID:     1207,
			Quote:  "Everything Has Its Limit - Iron Ore Cannot Be Educated Into Gold.",
			Author: "Mark Twain",
		},
		expectedResult: &api.Quote{
			ID:     1207,
			Quote:  "Everything Has Its Limit - Iron Ore Cannot Be Educated Into Gold.",
			Author: "Mark Twain",
		},
	}))

	t.Run("returns a server error when something went wrong", run(Test{
		mockedServiceError: errors.New("something went wrong"),
		expectedResult: &api.InternalServerErrror{
			Message: "unknown error",
		},
	}))
}
