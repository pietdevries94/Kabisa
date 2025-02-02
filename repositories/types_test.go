package repositories

import (
	"io"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockedHttpClient struct {
	mock.Mock
}

// Get is mocked here. Because it is rather annoying to check the full request that is passed to the mock,
// We use the url to match it.
func (m *MockedHttpClient) Do(req *http.Request) (resp *http.Response, err error) {
	url := req.URL.String()
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
