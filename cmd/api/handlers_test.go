package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/pietdevries94/Kabisa/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedQuoteService struct {
	mock.Mock
}

// Get is fully mocked here, returning with
func (m *MockedQuoteService) GetRandomQuote() (*models.Quote, error) {
	args := m.Called()
	return args.Get(0).(*models.Quote), args.Error(1)
}

func TestApplicationGetQuoteHandler(t *testing.T) {
	type Test struct {
		mockedServiceQuote *models.Quote
		mockedServiceError error
		expectedStatusCode int
		expectedBody       string
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

			// We build a request and pass it trough the handler using a recorder
			req, err := http.NewRequest("GET", "/quote", http.NoBody)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			http.HandlerFunc(app.GetQuoteHandler).ServeHTTP(rr, req)

			assrt := assert.New(t) // we rename this because we want to prevent shadowing

			assrt.Equal(tt.expectedStatusCode, rr.Code)
			assrt.Equal(tt.expectedBody, rr.Body.String())

			// We only expect the content type to be set if there is actual content
			var expectedContentType = "application/json"
			if rr.Body.String() == "" {
				expectedContentType = ""
			}
			assrt.Equal(expectedContentType, rr.Result().Header.Get("Content-Type"))
		}
	}

	t.Run("returns a quote", run(Test{
		mockedServiceQuote: &models.Quote{
			ID:     1207,
			Quote:  "Everything Has Its Limit - Iron Ore Cannot Be Educated Into Gold.",
			Author: "Mark Twain",
		},
		expectedBody:       `{"id":1207,"quote":"Everything Has Its Limit - Iron Ore Cannot Be Educated Into Gold.","author":"Mark Twain"}` + "\n",
		expectedStatusCode: http.StatusOK,
	}))

	t.Run("returns a server error when something went wrong", run(Test{
		mockedServiceError: errors.New("something went wrong"),
		expectedStatusCode: http.StatusInternalServerError,
	}))
}
