package repositories

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pietdevries94/Kabisa/database"
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
			// We get a new fresh inmem db for each test
			db := database.Init(&logger, ":memory:")

			res, err := NewQuoteGameRepo(&logger, db).CreateQuoteGame(context.TODO(), tt.quotes)

			assrt := assert.New(t) // we rename to prevent shadowing
			if tt.expectedError != nil {
				require.ErrorContains(t, err, tt.expectedError.Error())
				assrt.Equal(tt.expectedResult, res)
				return
			}

			require.NoError(t, err)

			// Because we can't predict a new uuid, and it's overkill to dependency inject the uuid function,
			// we check the uuid for validity and add it to the expectedResult
			assrt.NoError(uuid.Validate(res.ID.String()))
			tt.expectedResult.ID = res.ID

			assrt.Equal(tt.expectedResult, res)

			// Finally we check if the data is in the db in the expected way
			var id uuid.UUID
			var quote1_id, quote2_id, quote3_id int
			var ts time.Time
			err = db.QueryRow("select id, quote1_id, quote2_id, quote3_id, created_at from quote_game where id = ?", res.ID).
				Scan(&id, &quote1_id, &quote2_id, &quote3_id, &ts)
			require.NoError(t, err)

			assrt.Equal(res.ID, id)
			assrt.Equal(tt.expectedResult.Quotes[0].ID, quote1_id)
			assrt.Equal(tt.expectedResult.Quotes[1].ID, quote2_id)
			assrt.Equal(tt.expectedResult.Quotes[2].ID, quote3_id)
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

	t.Run("errors when not given 3 quotes", run(Test{
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
		},
		expectedError: errors.New("number of quotes should be 3. Given: 2"),
	}))
}

// TODO: implement test for ValidateIDAndAnswerIDs and ValidateAnswersAndCreateGameResult
