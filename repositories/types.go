package repositories

import "net/http"

// httpClient is an interface containing all the functions of http.Client that are used by the repositories
type httpClient interface {
	Get(url string) (resp *http.Response, err error)
}
