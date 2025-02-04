package repositories

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/pietdevries94/Kabisa/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDummyJsonRepo_GetRandomQuotes(t *testing.T) {
	type Test struct {
		amount              int
		mockedResponse      *http.Response
		mockedError         error
		expectedResult      []*models.Quote
		expectedError       error
		expectApiToBeCalled bool
	}

	run := func(tt Test) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()

			mockedHttpClient := new(MockedHttpClient)

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

			if tt.expectApiToBeCalled {
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
		expectApiToBeCalled: true,
	}))

	t.Run("returns an error when client.Get returns an error", run(Test{
		amount:              1,
		mockedError:         http.ErrHandlerTimeout,
		expectedError:       http.ErrHandlerTimeout,
		expectApiToBeCalled: true,
	}))

	t.Run("returns an error when the client.Get response returns no body", run(Test{
		amount:              1,
		mockedResponse:      CreateMockedResponse(http.StatusOK, nil),
		expectedError:       errors.New("no body received"),
		expectApiToBeCalled: true,
	}))

	t.Run("returns an error when the client.Get response doesn't return a 200", run(Test{
		amount:              1,
		mockedResponse:      CreateMockedResponse(http.StatusTeapot, bytes.NewBufferString("{}")),
		expectedError:       errors.New("unexpected status code received: 418"),
		expectApiToBeCalled: true,
	}))

	t.Run("returns an error when the cliet.Get response body is not valid json", run(Test{
		amount:              1,
		mockedResponse:      CreateMockedResponse(http.StatusOK, bytes.NewBufferString("<Quote>I'm XML<Quote>")),
		expectedError:       errors.New("unexpected error when decoding result to models.Quote"),
		expectApiToBeCalled: true,
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
		expectApiToBeCalled: true,
	}))

	t.Run("returns an error when asking for zero quotes", run(Test{
		amount:              0,
		expectedError:       errors.New("amount should be between 1 and 10. Given: 0"),
		expectApiToBeCalled: false,
	}))
	t.Run("returns an error when asking for negative quotes", run(Test{
		amount:              -5,
		expectedError:       errors.New("amount should be between 1 and 10. Given: -5"),
		expectApiToBeCalled: false,
	}))
	t.Run("returns an error when asking for over 10 quotes", run(Test{
		amount:              11,
		expectedError:       errors.New("amount should be between 1 and 10. Given: 11"),
		expectApiToBeCalled: false,
	}))
}

func TestDummyJsonRepo_GetQuote(t *testing.T) {
	type Test struct {
		id                    int
		mockedResponse        *http.Response
		mockedError           error
		expectedResult        *models.Quote
		expectedError         error
		expectErrorToBePublic bool
		expectApiToBeCalled   bool
	}

	run := func(tt Test) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()

			mockedHttpClient := new(MockedHttpClient)

			url := fmt.Sprintf("https://dummyjson.com/quotes/%d", tt.id)
			mockedHttpClient.On("Do", url).
				Once().
				Return(tt.mockedResponse, tt.mockedError)

			// We inject the mocked repo and expect to get the same quote back, but now as a struct
			logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
			res, err := NewDummyJsonRepo(&logger, mockedHttpClient).GetQuote(context.TODO(), tt.id)

			if tt.expectedError != nil {
				// We want to explicitly check if the error going out was meant to be a public type
				if tt.expectErrorToBePublic {
					assert.IsType(t, &models.PublicError{}, err)
				} else {
					assert.NotEqual(t, reflect.TypeOf(&models.PublicError{}), reflect.TypeOf(err))
				}
				require.ErrorContains(t, err, tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedResult, res)

			if tt.expectApiToBeCalled {
				mockedHttpClient.AssertCalled(t, "Do", url)
			} else {
				mockedHttpClient.AssertNotCalled(t, "Do", url)
			}
		}
	}

	t.Run("returns a quote when receiving expected response from api", run(Test{
		id:             414,
		mockedResponse: CreateMockedResponse(http.StatusOK, bytes.NewBufferString(`{"id":414,"quote":"When We Lose One Blessing, Another Is Often Most Unexpectedly Given In Its Place.","author":"C. S. Lewis"}`)),
		expectedResult: &models.Quote{
			ID:     414,
			Quote:  "When We Lose One Blessing, Another Is Often Most Unexpectedly Given In Its Place.",
			Author: "C. S. Lewis",
		},
		expectApiToBeCalled: true,
	}))

	t.Run("returns a public error when retrieving quote that doesn't exist", run(Test{
		id:                    414,
		mockedResponse:        CreateMockedResponse(http.StatusNotFound, bytes.NewBufferString(`{"message":"Quote with id '414' not found"}`)),
		expectedError:         models.NewPublicError("unknown_quote_id: 414"),
		expectErrorToBePublic: true,
		expectApiToBeCalled:   true,
	}))

	t.Run("returns an error when client.Get returns an error", run(Test{
		id:                  414,
		mockedError:         http.ErrHandlerTimeout,
		expectedError:       http.ErrHandlerTimeout,
		expectApiToBeCalled: true,
	}))

	t.Run("returns an error when the client.Get response returns no body", run(Test{
		id:                  414,
		mockedResponse:      CreateMockedResponse(http.StatusOK, nil),
		expectedError:       errors.New("no body received"),
		expectApiToBeCalled: true,
	}))

	t.Run("returns an error when the client.Get response doesn't return a 200", run(Test{
		id:                  414,
		mockedResponse:      CreateMockedResponse(http.StatusTeapot, bytes.NewBufferString("{}")),
		expectedError:       errors.New("unexpected status code received: 418"),
		expectApiToBeCalled: true,
	}))

	t.Run("returns an error when the cliet.Get response body is not valid json", run(Test{
		id:                  414,
		mockedResponse:      CreateMockedResponse(http.StatusOK, bytes.NewBufferString("<Quote>I'm XML<Quote>")),
		expectedError:       errors.New("unexpected error when decoding result to models.Quote"),
		expectApiToBeCalled: true,
	}))
}

// TODO: implement test for GetQuotes
func TestDummyJsonRepo_GetQuotes(t *testing.T) {
	type MockSets struct {
		resp *http.Response
		err  error
	}
	type Test struct {
		ids                   []int
		mockSets              map[int]MockSets
		expectedResult        map[int]*models.Quote
		expectedError         error
		expectErrorToBePublic bool
		expectApiToBeCalled   map[int]bool
	}

	run := func(tt Test) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()

			mockedHttpClient := new(MockedHttpClient)

			mockedUrls := map[int]string{}
			for id, set := range tt.mockSets {
				url := fmt.Sprintf("https://dummyjson.com/quotes/%d", id)
				mockedHttpClient.On("Do", url).
					Once().
					Return(set.resp, set.err)
				mockedUrls[id] = url
			}

			// We inject the mocked repo and expect to get the same quote back, but now as a struct
			logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
			res, err := NewDummyJsonRepo(&logger, mockedHttpClient).GetQuotes(context.TODO(), tt.ids)

			if tt.expectedError != nil {
				// We want to explicitly check if the error going out was meant to be a public type
				if tt.expectErrorToBePublic {
					assert.IsType(t, &models.PublicError{}, err)
				} else {
					assert.NotEqual(t, reflect.TypeOf(&models.PublicError{}), reflect.TypeOf(err))
				}
				require.ErrorContains(t, err, tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedResult, res)

			// We loop trough all mocked ids and check if the url got called
			for id, url := range mockedUrls {
				if tt.expectApiToBeCalled[id] {
					mockedHttpClient.AssertCalled(t, "Do", url)
				} else {
					mockedHttpClient.AssertNotCalled(t, "Do", url)
				}
			}
		}
	}

	t.Run("returns a quote when receiving expected response from api", run(Test{
		ids: []int{414, 172},
		mockSets: map[int]MockSets{
			414: {resp: CreateMockedResponse(http.StatusOK, bytes.NewBufferString(`{"id":414,"quote":"When We Lose One Blessing, Another Is Often Most Unexpectedly Given In Its Place.","author":"C. S. Lewis"}`))},
			172: {resp: CreateMockedResponse(http.StatusOK, bytes.NewBufferString(`{"id":172,"quote":"The only lasting beauty is the beauty of the heart.","author":"Rumi"}`))},
		},
		expectedResult: map[int]*models.Quote{
			414: {
				ID:     414,
				Quote:  "When We Lose One Blessing, Another Is Often Most Unexpectedly Given In Its Place.",
				Author: "C. S. Lewis",
			},
			172: {
				ID:     172,
				Quote:  "The only lasting beauty is the beauty of the heart.",
				Author: "Rumi",
			},
		},
		expectApiToBeCalled: map[int]bool{
			414: true,
			172: true,
		},
	}))

	t.Run("returns a public error when retrieving quote that doesn't exist", run(Test{
		ids: []int{414, 172},
		mockSets: map[int]MockSets{
			414: {resp: CreateMockedResponse(http.StatusNotFound, bytes.NewBufferString(`{"message":"Quote with id '414' not found"}`))},
			172: {resp: CreateMockedResponse(http.StatusOK, bytes.NewBufferString(`{"id":172,"quote":"The only lasting beauty is the beauty of the heart.","author":"Rumi"}`))},
		},
		expectedError:         models.NewPublicError("unknown_quote_id: 414"),
		expectErrorToBePublic: true,
		expectApiToBeCalled: map[int]bool{
			414: true,
			172: false,
		},
	}))

	t.Run("returns an error when client.Get returns an error", run(Test{
		ids: []int{414, 172},
		mockSets: map[int]MockSets{
			414: {err: http.ErrHandlerTimeout},
			172: {resp: CreateMockedResponse(http.StatusOK, bytes.NewBufferString(`{"id":172,"quote":"The only lasting beauty is the beauty of the heart.","author":"Rumi"}`))},
		},
		expectedError:         http.ErrHandlerTimeout,
		expectErrorToBePublic: false,
		expectApiToBeCalled: map[int]bool{
			414: true,
			172: false,
		},
	}))
}
