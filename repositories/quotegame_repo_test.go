package repositories

import (
	"context"
	"database/sql"
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
			defer db.Close()

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

func TestQuoteGameRepo_ValidateIDAndAnswerIDs(t *testing.T) {
	type Test struct {
		id             uuid.UUID
		answers        models.QuoteGameAnswerMap
		prepareDB      func(*sql.DB)
		expectedResult []int
		expectedError  error
	}

	run := func(tt Test) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()

			logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
			// We get a new fresh inmem db for each test
			db := database.Init(&logger, ":memory:")
			defer db.Close()
			// And seed it
			db.Exec( //nolint:errcheck // this is a test
				"insert into quote_game(id, quote1_id, quote2_id, quote3_id, created_at) values (?,?,?,?,?)",
				uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
				12,
				72,
				33,
				time.Now(),
			)
			if tt.prepareDB != nil {
				tt.prepareDB(db)
			}

			res, err := NewQuoteGameRepo(&logger, db).ValidateIDAndAnswerIDs(context.TODO(), tt.id, tt.answers)

			assrt := assert.New(t) // we rename to prevent shadowing
			if tt.expectedError != nil {
				require.ErrorContains(t, err, tt.expectedError.Error())
				assrt.Equal(tt.expectedResult, res)
				return
			}

			require.NoError(t, err)
			assrt.Equal(tt.expectedResult, res)
		}
	}

	t.Run("returns the quote ids in the order of the db when everything matches", run(Test{
		id: uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
		answers: models.QuoteGameAnswerMap{
			33: "Max",
			12: "Bob",
			72: "Jan",
		},
		expectedResult: []int{12, 72, 33},
	}))

	t.Run("throws error if the id doesn't exist", run(Test{
		id: uuid.MustParse("03f17f15-eeee-eeee-eeee-039f2f18373e"),
		answers: models.QuoteGameAnswerMap{
			33: "Max",
			12: "Bob",
			72: "Jan",
		},
		expectedError: models.ErrQuoteGameIdNotFound,
	}))

	t.Run("throws error if the answer ids don't match with the game", run(Test{
		id: uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
		answers: models.QuoteGameAnswerMap{
			99: "Max",
			12: "Bob",
			72: "Jan",
		},
		expectedError: models.ErrInvalidQuoteID,
	}))

	t.Run("throws error if the game is expired", run(Test{
		id: uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
		answers: models.QuoteGameAnswerMap{
			33: "Max",
			12: "Bob",
			72: "Jan",
		},
		prepareDB: func(db *sql.DB) {
			db.Exec( //nolint:errcheck // this is a test
				"update quote_game set created_at=? where id=?",
				time.Now().Add(-10*time.Minute),
				uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
			)
		},
		expectedError: models.ErrQuoteGameIdNotFound,
	}))

	t.Run("throws error if the game is already completed", run(Test{
		id: uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
		answers: models.QuoteGameAnswerMap{
			33: "Max",
			12: "Bob",
			72: "Jan",
		},
		prepareDB: func(db *sql.DB) {
			db.Exec( //nolint:errcheck // this is a test
				"update quote_game set completed_at=? where id=?",
				time.Now(),
				uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
			)
		},
		expectedError: models.ErrQuoteGameIdNotFound,
	}))
}

func TestQuoteGameRepo_ValidateAnswersAndCreateGameResult(t *testing.T) {
	type Test struct {
		id             uuid.UUID
		quoteIDs       []int
		quotes         map[int]*models.Quote
		answers        models.QuoteGameAnswerMap
		expectedResult *models.QuoteGameResult
		expectedError  error
	}

	run := func(tt Test) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()

			logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
			// We get a new fresh inmem db for each test
			db := database.Init(&logger, ":memory:")
			defer db.Close()
			// And seed it
			db.Exec( //nolint:errcheck // this is a test
				"insert into quote_game(id, quote1_id, quote2_id, quote3_id, created_at) values (?,?,?,?,?)",
				uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
				12,
				72,
				33,
				time.Now(),
			)

			res, err := NewQuoteGameRepo(&logger, db).ValidateAnswersAndCreateGameResult(context.TODO(), tt.id, tt.quoteIDs, tt.quotes, tt.answers)

			assrt := assert.New(t) // we rename to prevent shadowing
			req := require.New(t)
			if tt.expectedError != nil {
				req.ErrorContains(err, tt.expectedError.Error())
				assrt.Equal(tt.expectedResult, res)
				return
			}

			req.NoError(err)
			assrt.Equal(tt.expectedResult, res)

			// We want to check if the state is actually set in the db
			var q1_correct, q2_correct, q3_correct sql.NullBool
			var completed_at sql.NullTime
			err = db.QueryRow("select quote1_correct, quote2_correct, quote3_correct, completed_at from quote_game where id = ?", tt.id).
				Scan(&q1_correct, &q2_correct, &q3_correct, &completed_at)
			req.NoError(err)

			req.True(q1_correct.Valid)
			req.True(q2_correct.Valid)
			req.True(q3_correct.Valid)
			req.True(completed_at.Valid)

			assrt.Equal(tt.expectedResult.Answers[0].Correct, q1_correct.Bool)
			assrt.Equal(tt.expectedResult.Answers[1].Correct, q2_correct.Bool)
			assrt.Equal(tt.expectedResult.Answers[2].Correct, q3_correct.Bool)
		}
	}

	t.Run("checks the answers and stores them in the database", run(Test{
		id: uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
		answers: models.QuoteGameAnswerMap{
			33: "Max",
			12: "Bob",
			72: "Jan",
		},
		quoteIDs: []int{12, 72, 33},
		quotes: map[int]*models.Quote{
			33: {Author: "Max", Quote: "Hey", ID: 33},
			12: {Author: "Bob", Quote: "Hi", ID: 12},
			72: {Author: "Someone else", Quote: "Bye", ID: 72},
		},
		expectedResult: &models.QuoteGameResult{
			ID: uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
			Answers: []*models.QuoteGameActualAnswer{
				{Quote: models.Quote{Author: "Bob", Quote: "Hi", ID: 12}, Correct: true},
				{Quote: models.Quote{Author: "Someone else", Quote: "Bye", ID: 72}, Correct: false},
				{Quote: models.Quote{Author: "Max", Quote: "Hey", ID: 33}, Correct: true},
			},
		},
	}))
}
