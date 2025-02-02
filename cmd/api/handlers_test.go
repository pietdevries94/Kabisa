package main

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/pietdevries94/Kabisa/models"
	"github.com/pietdevries94/Kabisa/openapi"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplication_GetRandomQuote(t *testing.T) {
	type Test struct {
		mockedServiceQuote *models.Quote
		mockedServiceError error
		expectedResult     openapi.GetRandomQuoteRes
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
		expectedResult: &openapi.Quote{
			ID:     1207,
			Quote:  "Everything Has Its Limit - Iron Ore Cannot Be Educated Into Gold.",
			Author: "Mark Twain",
		},
	}))

	t.Run("returns a server error when something went wrong", run(Test{
		mockedServiceError: errors.New("something went wrong"),
		expectedResult: &openapi.R500{
			Message: "unknown_error",
		},
	}))
}

func TestApplication_CreateQuoteGame(t *testing.T) {
	type Test struct {
		mockedServiceQuote *models.QuoteGame
		mockedServiceError error
		expectedResult     openapi.CreateNewQuoteGameRes
	}

	run := func(tt Test) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()

			// first we bootstrap a minimal version of the application, needed for the handler
			mockedQuoteService := new(MockedQuoteService)
			mockedQuoteService.On("CreateQuoteGame").Once().Return(tt.mockedServiceQuote, tt.mockedServiceError)

			logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
			app := application{
				logger:       &logger,
				quoteService: mockedQuoteService,
			}

			// We now run the handler and validate the result
			res, err := app.CreateNewQuoteGame(context.TODO())
			require.NoError(t, err)
			assert.Equal(t, tt.expectedResult, res)
		}
	}

	t.Run("returns a quote game", run(Test{
		mockedServiceQuote: &models.QuoteGame{
			ID: uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
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
		expectedResult: &openapi.CreateNewQuoteGameOK{
			ID: "03f17f15-5d0a-49ea-aa05-039f2f18373e",
			Quotes: []openapi.QuoteWithoutAuthor{
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

	t.Run("returns a server error when something went wrong", run(Test{
		mockedServiceError: errors.New("something went wrong"),
		expectedResult: &openapi.R500{
			Message: "unknown_error",
		},
	}))
}

func TestApplication_SubmitAnswerForQuoteGame(t *testing.T) {
	type Test struct {
		answers                           []openapi.QuoteGameAnswer
		params                            openapi.SubmitAnswerForQuoteGameParams
		expectedMockedServiceInputID      uuid.UUID
		expectedMockedServiceInputAnswers models.QuoteGameAnswerMap
		mockedServiceResult               *models.QuoteGameResult
		mockedServiceError                error
		expectedResult                    openapi.SubmitAnswerForQuoteGameRes
	}

	run := func(tt Test) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()

			// first we bootstrap a minimal version of the application, needed for the handler
			mockedQuoteService := new(MockedQuoteService)
			mockedQuoteService.On("SubmitAnswerToQuoteGame", tt.expectedMockedServiceInputID, tt.expectedMockedServiceInputAnswers).
				Once().
				Return(tt.mockedServiceResult, tt.mockedServiceError)

			logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
			app := application{
				logger:       &logger,
				quoteService: mockedQuoteService,
			}

			// We now run the handler and validate the result
			res, err := app.SubmitAnswerForQuoteGame(context.TODO(), tt.answers, tt.params)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedResult, res)
		}
	}

	t.Run("returns a quote game", run(Test{
		answers: []openapi.QuoteGameAnswer{
			{ID: 54, Author: "A name"},
			{ID: 43, Author: "A different name"},
			{ID: 2, Author: "Bob"},
		},
		params: openapi.SubmitAnswerForQuoteGameParams{
			ID: "03f17f15-5d0a-49ea-aa05-039f2f18373e",
		},
		expectedMockedServiceInputID: uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
		expectedMockedServiceInputAnswers: models.QuoteGameAnswerMap{
			54: "A name",
			43: "A different name",
			2:  "Bob",
		},
		mockedServiceResult: &models.QuoteGameResult{
			ID: uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
			Answers: []*models.QuoteGameActualAnswer{
				{Quote: models.Quote{ID: 54, Author: "George", Quote: "Hello!"}, Correct: false},
				{Quote: models.Quote{ID: 43, Author: "William", Quote: "Hi!"}, Correct: false},
				{Quote: models.Quote{ID: 2, Author: "Bob", Quote: "Bye!"}, Correct: true},
			},
		},
		expectedResult: &openapi.QuoteGameResult{
			ID: "03f17f15-5d0a-49ea-aa05-039f2f18373e",
			Answers: []openapi.QuoteGameResultAnswersItem{
				{ID: 54, Correct: false, ActualAuthor: "George"},
				{ID: 43, Correct: false, ActualAuthor: "William"},
				{ID: 2, Correct: true, ActualAuthor: "Bob"},
			},
		},
	}))

	t.Run("returns an 500 if the service errors without a public error", run(Test{
		answers: []openapi.QuoteGameAnswer{
			{ID: 54, Author: "A name"},
			{ID: 43, Author: "A different name"},
			{ID: 2, Author: "Bob"},
		},
		params: openapi.SubmitAnswerForQuoteGameParams{
			ID: "03f17f15-5d0a-49ea-aa05-039f2f18373e",
		},
		expectedMockedServiceInputID: uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
		expectedMockedServiceInputAnswers: models.QuoteGameAnswerMap{
			54: "A name",
			43: "A different name",
			2:  "Bob",
		},
		mockedServiceError: errors.New("a crazy error"),
		expectedResult: &openapi.R500{
			Message: "unknown_error",
		},
	}))

	t.Run("returns a 404 if the id is not parseable as a uuid v4", run(Test{
		answers: []openapi.QuoteGameAnswer{
			{ID: 54, Author: "A name"},
			{ID: 43, Author: "A different name"},
			{ID: 2, Author: "Bob"},
		},
		params: openapi.SubmitAnswerForQuoteGameParams{
			ID: "nope",
		},
		expectedResult: &openapi.R404{
			Message: "not_found",
		},
	}))

	t.Run("returns a 404 if the error returned by the service is a quote_game_id_not_found error", run(Test{
		answers: []openapi.QuoteGameAnswer{
			{ID: 54, Author: "A name"},
			{ID: 43, Author: "A different name"},
			{ID: 2, Author: "Bob"},
		},
		params: openapi.SubmitAnswerForQuoteGameParams{
			ID: "03f17f15-5d0a-49ea-aa05-039f2f18373e",
		},
		expectedMockedServiceInputID: uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
		expectedMockedServiceInputAnswers: models.QuoteGameAnswerMap{
			54: "A name",
			43: "A different name",
			2:  "Bob",
		},
		mockedServiceError: models.ErrQuoteGameIdNotFound,
		expectedResult: &openapi.R404{
			Message: "not_found",
		},
	}))

	t.Run("returns a 422 if any other public error is given", run(Test{
		answers: []openapi.QuoteGameAnswer{
			{ID: 54, Author: "A name"},
			{ID: 43, Author: "A different name"},
			{ID: 2, Author: "Bob"},
		},
		params: openapi.SubmitAnswerForQuoteGameParams{
			ID: "03f17f15-5d0a-49ea-aa05-039f2f18373e",
		},
		expectedMockedServiceInputID: uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
		expectedMockedServiceInputAnswers: models.QuoteGameAnswerMap{
			54: "A name",
			43: "A different name",
			2:  "Bob",
		},
		mockedServiceError: models.ErrInvalidQuoteID,
		expectedResult: &openapi.R422{
			Message: "invalid_quote_id",
		},
	}))
}
