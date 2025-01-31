package services

import (
	"errors"
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
	type Test struct {
		mockedJsonRepoQuote *models.Quote
		mockedJsonRepoError error
		expectedResult      *models.Quote
		expectedError       error
	}
	run := func(tt Test) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()

			mockedDummyJsonRepo := new(MockedDummyJsonRepo)
			mockedDummyJsonRepo.On("GetRandomQuote").Once().Return(tt.mockedJsonRepoQuote, tt.mockedJsonRepoError)

			// We inject the mocked repo into the service and expect the same quote back
			logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
			res, err := NewQuoteService(&logger, mockedDummyJsonRepo).GetRandomQuote()

			if tt.expectedError != nil {
				require.ErrorContains(t, err, tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedResult, res)
		}
	}

	t.Run("returns quote from dummyJsonRepo", run(Test{
		// The mocked quote is based on an actual response from the underlying api
		mockedJsonRepoQuote: &models.Quote{
			ID:     663,
			Quote:  "Never Mistake Motion For Action.",
			Author: "Ernest Hemingway",
		},
		expectedResult: &models.Quote{
			ID:     663,
			Quote:  "Never Mistake Motion For Action.",
			Author: "Ernest Hemingway",
		},
	}))

	t.Run("passes trough an error from dummyJsonRepo", run(
		Test{
			mockedJsonRepoError: errors.New("this is an error"),
			expectedError:       errors.New("this is an error"),
		},
	))
}
