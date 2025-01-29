package repositories

import (
	"encoding/json"
	"net/http"

	"github.com/pietdevries94/Kabisa/models"
)

// DummyJsonRepo handles all calls to the dummyjson.com api.
type DummyJsonRepo interface {
	// GetRandomQuote retrieves a quote using the random feature of dummyjson
	GetRandomQuote() (*models.Quote, error)
}

// httpClient is an interface containing all the functions of http.Client that are used by the repositories
type httpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type dummyJsonRepo struct {
	client httpClient
}

// NewDummyJsonRepo returns a new DummyJsonRepo, which handles all calls to the dummyjson.com api.
func NewDummyJsonRepo(client httpClient) DummyJsonRepo {
	return &dummyJsonRepo{client: client}
}

// GetRandomQuote retrieves a quote using the random feature of dummyjson
func (repo *dummyJsonRepo) GetRandomQuote() (*models.Quote, error) {
	resp, err := repo.client.Get("https://dummyjson.com/quotes/random")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// TODO: status check

	var quote *models.Quote
	err = json.NewDecoder(resp.Body).Decode(&quote)
	if err != nil {
		return nil, err
	}

	return quote, nil
}
