package repositories

import "net/http"

// httpClient is an interface containing all the functions of http.Client that are used by the repositories
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
