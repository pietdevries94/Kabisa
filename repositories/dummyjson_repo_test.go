package repositories

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/pietdevries94/Kabisa/models"
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
	t.Run("returns quote when receiving expected response from api", func(t *testing.T) {
		mockedClient := new(MockedHttpClient)

		// The body is a copy of an actual response from the api
		mockedResponse := CreateMockedResponse(http.StatusOK, bytes.NewBufferString(`{"id":414,"quote":"When We Lose One Blessing, Another Is Often Most Unexpectedly Given In Its Place.","author":"C. S. Lewis"}`))
		mockedClient.On("Get", "https://dummyjson.com/quotes/random").
			Once().
			Return(mockedResponse, nil)

		// We inject the mocked repo and expect to get the same quote back, but now as a struct
		res, err := NewDummyJsonRepo(mockedClient).GetRandomQuote()
		require.NoError(t, err)
		assert.Equal(t, &models.Quote{
			ID:     414,
			Quote:  "When We Lose One Blessing, Another Is Often Most Unexpectedly Given In Its Place.",
			Author: "C. S. Lewis",
		}, res)
	})
	// TODO extend tests and convert to closure tests
}
