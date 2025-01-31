package repositories

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/pietdevries94/Kabisa/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockedHttpClient struct {
	mock.Mock
}

// Get is fully mocked here, returning with
func (m *MockedHttpClient) Get(url string) (resp *http.Response, err error) {
	args := m.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}

func CreateMockedResponse(statusCode int, bodyReader io.Reader) *http.Response {
	resp := &http.Response{
		Status:     http.StatusText(statusCode),
		StatusCode: statusCode,
	}
	// We only want to set the body if there actually is given one
	if bodyReader != nil {
		resp.Body = io.NopCloser(bodyReader)
	}
	return resp
}

func TestDummyJsonRepoGetRandomQuote(t *testing.T) {
	type Test struct {
		mockedResponse *http.Response
		mockedError    error
		expectedResult *models.Quote
		expectedError  error
	}

	run := func(tt Test) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()

			mockedHttpClient := new(MockedHttpClient)

			// The body is a copy of an actual response from the api
			mockedHttpClient.On("Get", "https://dummyjson.com/quotes/random").
				Once().
				Return(tt.mockedResponse, tt.mockedError)

			// We inject the mocked repo and expect to get the same quote back, but now as a struct
			logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
			res, err := NewDummyJsonRepo(&logger, mockedHttpClient).GetRandomQuote()

			if tt.expectedError != nil {
				require.ErrorContains(t, err, tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedResult, res)
		}
	}

	t.Run("returns quote when receiving expected response from api", run(Test{
		mockedResponse: CreateMockedResponse(http.StatusOK, bytes.NewBufferString(`{"id":414,"quote":"When We Lose One Blessing, Another Is Often Most Unexpectedly Given In Its Place.","author":"C. S. Lewis"}`)),
		expectedResult: &models.Quote{
			ID:     414,
			Quote:  "When We Lose One Blessing, Another Is Often Most Unexpectedly Given In Its Place.",
			Author: "C. S. Lewis",
		},
	}))

	t.Run("returns an error when client.Get returns an error", run(Test{
		mockedError:   http.ErrHandlerTimeout,
		expectedError: http.ErrHandlerTimeout,
	}))

	t.Run("returns an error when the client.Get response returns no body", run(Test{
		mockedResponse: CreateMockedResponse(http.StatusOK, nil),
		expectedError:  errors.New("no body received"),
	}))

	t.Run("returns an error when the client.Get response doesn't return a 200", run(Test{
		mockedResponse: CreateMockedResponse(http.StatusTeapot, bytes.NewBufferString("{}")),
		expectedError:  errors.New("unexpected status code received: 418"),
	}))

	t.Run("returns an error when the cliet.Get response body is not valid json", run(Test{
		mockedResponse: CreateMockedResponse(http.StatusOK, bytes.NewBufferString("<Quote>I'm XML<Quote>")),
		expectedError:  errors.New("unexpected error when decoding result to models.Quote"),
	}))
}
