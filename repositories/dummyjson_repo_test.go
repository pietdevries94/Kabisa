package repositories

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/pietdevries94/Kabisa/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDummyJsonRepo_GetRandomQuotes(t *testing.T) {
	type Test struct {
		amount                int
		mockedResponse        *http.Response
		mockedError           error
		expectedResult        []*models.Quote
		expectedError         error
		expectedApiToBeCalled bool
	}

	run := func(tt Test) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()

			mockedHttpClient := new(MockedHttpClient)

			// The body is a copy of an actual response from the api
			url := fmt.Sprintf("https://dummyjson.com/quotes/random/%d", tt.amount)
			mockedHttpClient.On("Do", url).
				Once().
				Return(tt.mockedResponse, tt.mockedError)

			// We inject the mocked repo and expect to get the same quote back, but now as a struct
			logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
			res, err := NewDummyJsonRepo(&logger, mockedHttpClient).GetRandomQuotes(context.TODO(), tt.amount)

			if tt.expectedError != nil {
				require.ErrorContains(t, err, tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedResult, res)

			if tt.expectedApiToBeCalled {
				mockedHttpClient.AssertCalled(t, "Do", url)
			} else {
				mockedHttpClient.AssertNotCalled(t, "Do", url)
			}
		}
	}

	t.Run("returns a quote when receiving expected response from api", run(Test{
		amount:         1,
		mockedResponse: CreateMockedResponse(http.StatusOK, bytes.NewBufferString(`[{"id":414,"quote":"When We Lose One Blessing, Another Is Often Most Unexpectedly Given In Its Place.","author":"C. S. Lewis"}]`)),
		expectedResult: []*models.Quote{
			{
				ID:     414,
				Quote:  "When We Lose One Blessing, Another Is Often Most Unexpectedly Given In Its Place.",
				Author: "C. S. Lewis",
			},
		},
		expectedApiToBeCalled: true,
	}))

	t.Run("returns an error when client.Get returns an error", run(Test{
		amount:                1,
		mockedError:           http.ErrHandlerTimeout,
		expectedError:         http.ErrHandlerTimeout,
		expectedApiToBeCalled: true,
	}))

	t.Run("returns an error when the client.Get response returns no body", run(Test{
		amount:                1,
		mockedResponse:        CreateMockedResponse(http.StatusOK, nil),
		expectedError:         errors.New("no body received"),
		expectedApiToBeCalled: true,
	}))

	t.Run("returns an error when the client.Get response doesn't return a 200", run(Test{
		amount:                1,
		mockedResponse:        CreateMockedResponse(http.StatusTeapot, bytes.NewBufferString("{}")),
		expectedError:         errors.New("unexpected status code received: 418"),
		expectedApiToBeCalled: true,
	}))

	t.Run("returns an error when the cliet.Get response body is not valid json", run(Test{
		amount:                1,
		mockedResponse:        CreateMockedResponse(http.StatusOK, bytes.NewBufferString("<Quote>I'm XML<Quote>")),
		expectedError:         errors.New("unexpected error when decoding result to models.Quote"),
		expectedApiToBeCalled: true,
	}))

	t.Run("returns multiple quotes, when requested", run(Test{
		amount:         3,
		mockedResponse: CreateMockedResponse(http.StatusOK, bytes.NewBufferString(`[{"id":1386,"quote":"It Is Most Pleasant To Commit A Just Action Which Is Disagreeable To Someone Whom One Does Not Like.","author":"Victor Hugo"},{"id":172,"quote":"The only lasting beauty is the beauty of the heart.","author":"Rumi"},{"id":454,"quote":"Risk Comes From Not Knowing What You'Re Doing.","author":"Warren Buffett"}]`)),
		expectedResult: []*models.Quote{
			{
				ID:     1386,
				Quote:  "It Is Most Pleasant To Commit A Just Action Which Is Disagreeable To Someone Whom One Does Not Like.",
				Author: "Victor Hugo",
			},
			{
				ID:     172,
				Quote:  "The only lasting beauty is the beauty of the heart.",
				Author: "Rumi",
			},
			{
				ID:     454,
				Quote:  "Risk Comes From Not Knowing What You'Re Doing.",
				Author: "Warren Buffett",
			},
		},
		expectedApiToBeCalled: true,
	}))

	t.Run("returns an error when asking for zero quotes", run(Test{
		amount:                0,
		expectedError:         errors.New("amount should be between 1 and 10. Given: 0"),
		expectedApiToBeCalled: false,
	}))
	t.Run("returns an error when asking for negative quotes", run(Test{
		amount:                -5,
		expectedError:         errors.New("amount should be between 1 and 10. Given: -5"),
		expectedApiToBeCalled: false,
	}))
	t.Run("returns an error when asking for over 10 quotes", run(Test{
		amount:                11,
		expectedError:         errors.New("amount should be between 1 and 10. Given: 11"),
		expectedApiToBeCalled: false,
	}))
}

// TODO: implement test for GetQuote and GetQuotes
