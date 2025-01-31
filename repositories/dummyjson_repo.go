package repositories

import (
	"encoding/json"
	"net/http"

	"github.com/pietdevries94/Kabisa/models"
	"github.com/rs/zerolog"
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
	logger     *zerolog.Logger
	httpClient httpClient
}

// NewDummyJsonRepo returns a new DummyJsonRepo, which handles all calls to the dummyjson.com api.
func NewDummyJsonRepo(logger *zerolog.Logger, httpClient httpClient) DummyJsonRepo {
	return &dummyJsonRepo{
		logger:     logger,
		httpClient: httpClient,
	}
}

// GetRandomQuote retrieves a quote using the random feature of dummyjson
func (repo *dummyJsonRepo) GetRandomQuote() (*models.Quote, error) {
	resp, err := repo.httpClient.Get("https://dummyjson.com/quotes/random")
	if err != nil {
		repo.logger.Error().Err(err).Msg("unexpected error when retrieving random quote from api")
		return nil, err
	}
	defer resp.Body.Close()

	// TODO: status check

	var quote *models.Quote
	err = json.NewDecoder(resp.Body).Decode(&quote)
	if err != nil {
		repo.logger.Error().Err(err).Msg("unexpected error when decoding result to models.Quote")
		return nil, err
	}

	return quote, nil
}
