package repositories

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/pietdevries94/Kabisa/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuoteGameRepo_GetRandomQuotes(t *testing.T) {
	type Test struct {
		quotes         []*models.Quote
		expectedResult *models.QuoteGame
		expectedError  error
	}

	run := func(tt Test) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()
			logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
			res, err := NewQuoteGameRepo(&logger).CreateQuoteGame(tt.quotes)

			if tt.expectedError != nil {
				require.ErrorContains(t, err, tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}

			// Because we can't predict a new uuid, and it's overkill to dependency inject the uuid function,
			// we check the uuid for validity and add it to the expectedResult
			assert.NoError(t, uuid.Validate(res.ID.String()))
			tt.expectedResult.ID = res.ID

			assert.Equal(t, tt.expectedResult, res)
		}
	}

	t.Run("returns a quote when receiving expected response from api", run(Test{
		quotes: []*models.Quote{
			{
				ID:     70,
				Quote:  "The cure for pain is in the pain.",
				Author: "Rumi",
			},
			{
				ID:     905,
				Quote:  "Try as much as you can to mention death. For if you were having hard times in your life, then it would give you more hope and would ease things for you. And if you were having abundant affluence of living in luxury, then it would make it less luxurious.",
				Author: "Umar ibn Al-Khattāb (R.A)",
			},
			{
				ID:     451,
				Quote:  "We should not give up and we should not allow the problem to defeat us.",
				Author: "Abdul Kalam",
			},
		},
		expectedResult: &models.QuoteGame{
			Quotes: []*models.QuoteWithoutAuthor{
				{
					ID:    70,
					Quote: "The cure for pain is in the pain.",
				},
				{
					ID:    905,
					Quote: "Try as much as you can to mention death. For if you were having hard times in your life, then it would give you more hope and would ease things for you. And if you were having abundant affluence of living in luxury, then it would make it less luxurious.",
				},
				{
					ID:    451,
					Quote: "We should not give up and we should not allow the problem to defeat us.",
				},
			},
			Authors: []string{
				"Abdul Kalam",
				"Rumi",
				"Umar ibn Al-Khattāb (R.A)",
			},
		},
	}))
}
