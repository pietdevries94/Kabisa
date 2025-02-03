package services

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/pietdevries94/Kabisa/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuoteService_GetRandomQuote(t *testing.T) {
	type Test struct {
		mockedJsonRepoQuote []*models.Quote
		mockedJsonRepoError error
		expectedResult      *models.Quote
		expectedError       error
	}
	run := func(tt Test) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()

			mockedDummyJsonRepo := new(MockedDummyJsonRepo)
			mockedDummyJsonRepo.On("GetRandomQuotes", 1).
				Once().
				Return(tt.mockedJsonRepoQuote, tt.mockedJsonRepoError)

			// We inject the mocked repo into the service and expect the same quote back
			logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
			res, err := NewQuoteService(&logger, mockedDummyJsonRepo, nil).GetRandomQuote(context.TODO())

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
		mockedJsonRepoQuote: []*models.Quote{
			{
				ID:     663,
				Quote:  "Never Mistake Motion For Action.",
				Author: "Ernest Hemingway",
			},
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

	t.Run("throws an error if the dummyJsonRepo returns no quotes, nor an error", run(
		Test{
			expectedError: errors.New("dummyJsonRepo returned no quotes and no error"),
		},
	))
}

func TestQuoteService_CreateQuoteGame(t *testing.T) {
	type Test struct {
		mockedJsonRepoQuote  []*models.Quote
		mockedJsonRepoError  error
		mockedQuoteGame      *models.QuoteGame
		mockedQuoteGameError error
		expectedResult       *models.QuoteGame
		expectedError        error
	}
	run := func(tt Test) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()

			mockedDummyJsonRepo := new(MockedDummyJsonRepo)
			mockedDummyJsonRepo.On("GetRandomQuotes", 3).
				Once().
				Return(tt.mockedJsonRepoQuote, tt.mockedJsonRepoError)

			mockedQuoteGameRepo := new(MockedQuoteGameRepo)
			mockedQuoteGameRepo.On("CreateQuoteGame", tt.mockedJsonRepoQuote).
				Once().
				Return(tt.mockedQuoteGame, tt.mockedQuoteGameError)

			// We inject the mocked repos into the service and expect the same quote back
			logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
			res, err := NewQuoteService(&logger, mockedDummyJsonRepo, mockedQuoteGameRepo).CreateQuoteGame(context.TODO())

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
		mockedJsonRepoQuote: []*models.Quote{
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
		mockedQuoteGame: &models.QuoteGame{
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
		expectedResult: &models.QuoteGame{
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
	}))

	t.Run("passes trough an error from dummyJsonRepo", run(
		Test{
			mockedJsonRepoError: errors.New("this is an error"),
			expectedError:       errors.New("this is an error"),
		},
	))

	t.Run("passes trough an error from quoteGameRepo", run(
		Test{
			mockedJsonRepoQuote: []*models.Quote{
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
			mockedQuoteGameError: errors.New("this is an error"),
			expectedError:        errors.New("this is an error"),
		},
	))
}

func TestQuoteService_SubmitAnswerToQuoteGame(t *testing.T) {
	type Test struct {
		id                                       uuid.UUID
		answers                                  models.QuoteGameAnswerMap
		mockedValidateIDAndAnswerIDsResult       []int
		mockedValidateIDAndAnswerIDsError        error
		mockedGetQuotesResult                    map[int]*models.Quote
		mockedGetQuotesError                     error
		mockedValidateAnswersAndCreateGameResult *models.QuoteGameResult
		mockedValidateAnswersAndCreateGameError  error
		expectedResult                           *models.QuoteGameResult
		expectedError                            error
	}
	run := func(tt Test) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()

			mockedQuoteGameRepo := new(MockedQuoteGameRepo)
			mockedDummyJsonRepo := new(MockedDummyJsonRepo)

			mockedQuoteGameRepo.On("ValidateIDAndAnswerIDs", tt.id, tt.answers).
				Once().
				Return(tt.mockedValidateIDAndAnswerIDsResult, tt.mockedValidateIDAndAnswerIDsError)

			mockedDummyJsonRepo.On("GetQuotes", tt.mockedValidateIDAndAnswerIDsResult).
				Once().
				Return(tt.mockedGetQuotesResult, tt.mockedGetQuotesError)

			mockedQuoteGameRepo.On("ValidateAnswersAndCreateGameResult", tt.id, tt.mockedGetQuotesResult, tt.answers).
				Once().
				Return(tt.mockedValidateAnswersAndCreateGameResult, tt.mockedValidateAnswersAndCreateGameError)

			// We inject the mocked repos into the service and expect the same quote back
			logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
			res, err := NewQuoteService(&logger, mockedDummyJsonRepo, mockedQuoteGameRepo).
				SubmitAnswerToQuoteGame(context.TODO(), tt.id, tt.answers)

			if tt.expectedError != nil {
				require.ErrorContains(t, err, tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedResult, res)
		}
	}

	t.Run("returns the result, when given valid parameters", run(Test{
		id: uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
		answers: models.QuoteGameAnswerMap{
			54: "A name",
			43: "A different name",
			2:  "Bob",
		},
		mockedValidateIDAndAnswerIDsResult: []int{54, 43, 2},
		mockedGetQuotesResult: map[int]*models.Quote{
			54: {ID: 54, Author: "George", Quote: "Hello!"},
			43: {ID: 43, Author: "William", Quote: "Hi!"},
			2:  {ID: 2, Author: "Bob", Quote: "Bye!"},
		},
		mockedValidateAnswersAndCreateGameResult: &models.QuoteGameResult{
			ID: uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
			Answers: []*models.QuoteGameActualAnswer{
				{Quote: models.Quote{ID: 54, Author: "George", Quote: "Hello!"}, Correct: false},
				{Quote: models.Quote{ID: 43, Author: "William", Quote: "Hi!"}, Correct: false},
				{Quote: models.Quote{ID: 2, Author: "Bob", Quote: "Bye!"}, Correct: true},
			},
		},
		expectedResult: &models.QuoteGameResult{
			ID: uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
			Answers: []*models.QuoteGameActualAnswer{
				{Quote: models.Quote{ID: 54, Author: "George", Quote: "Hello!"}, Correct: false},
				{Quote: models.Quote{ID: 43, Author: "William", Quote: "Hi!"}, Correct: false},
				{Quote: models.Quote{ID: 2, Author: "Bob", Quote: "Bye!"}, Correct: true},
			},
		},
	}))

	t.Run("returns the error when ValidateAnswersAndCreateGameResult fails", run(Test{
		id: uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
		answers: models.QuoteGameAnswerMap{
			54: "A name",
			43: "A different name",
			2:  "Bob",
		},
		mockedValidateIDAndAnswerIDsResult: []int{54, 43, 2},
		mockedGetQuotesResult: map[int]*models.Quote{
			54: {ID: 54, Author: "George", Quote: "Hello!"},
			43: {ID: 43, Author: "William", Quote: "Hi!"},
			2:  {ID: 2, Author: "Bob", Quote: "Bye!"},
		},
		mockedValidateAnswersAndCreateGameError: errors.New("a brand new error"),
		expectedError:                           errors.New("a brand new error"),
	}))

	t.Run("returns the error when GetQuotes fails", run(Test{
		id: uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
		answers: models.QuoteGameAnswerMap{
			54: "A name",
			43: "A different name",
			2:  "Bob",
		},
		mockedValidateIDAndAnswerIDsResult: []int{54, 43, 2},
		mockedGetQuotesError:               errors.New("a brand new error"),
		expectedError:                      errors.New("a brand new error"),
	}))

	t.Run("returns the error when alidateIDAndAnswerIDs fails", run(Test{
		id: uuid.MustParse("03f17f15-5d0a-49ea-aa05-039f2f18373e"),
		answers: models.QuoteGameAnswerMap{
			54: "A name",
			43: "A different name",
			2:  "Bob",
		},
		mockedValidateIDAndAnswerIDsError: errors.New("a brand new error"),
		expectedError:                     errors.New("a brand new error"),
	}))
}
