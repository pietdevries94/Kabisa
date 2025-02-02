package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/pietdevries94/Kabisa/models"
	"github.com/rs/zerolog"
)

type DummyJsonRepo struct {
	logger     *zerolog.Logger
	httpClient httpClient
}

// NewDummyJsonRepo returns a new DummyJsonRepo, which handles all calls to the dummyjson.com api.
func NewDummyJsonRepo(logger *zerolog.Logger, httpClient httpClient) *DummyJsonRepo {
	return &DummyJsonRepo{
		logger:     logger,
		httpClient: httpClient,
	}
}

// GetRandomQuote retrieves a quote using the random feature of dummyjson
func (repo *DummyJsonRepo) GetRandomQuotes(_ context.Context, amount int) ([]*models.Quote, error) {
	if amount < 1 || amount > 10 {
		// The api only accepts 1 to 10
		return nil, fmt.Errorf("amount should be between 1 and 10. Given: %d", amount)
	}

	url := fmt.Sprintf("https://dummyjson.com/quotes/random/%d", amount)

	// TODO use the request
	resp, err := repo.httpClient.Get(url)
	if err != nil {
		repo.logger.Error().Err(err).Msg("unexpected error when retrieving random quote from api")
		return nil, errors.Join(errors.New("unexpected error when retrieving random quote from api"), err)
	}

	if resp.Body == nil {
		repo.logger.Error().Err(err).Int("status code", resp.StatusCode).Msg("no body received")
		return nil, errors.New("no body received")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		repo.logger.Error().Int("status code", resp.StatusCode).Msg("unexpected status code received")
		return nil, fmt.Errorf("unexpected status code received: %d", resp.StatusCode)
	}

	var quotes []*models.Quote
	err = json.NewDecoder(resp.Body).Decode(&quotes)
	if err != nil {
		repo.logger.Error().Err(err).Msg("unexpected error when decoding result to models.Quote")
		return nil, errors.Join(errors.New("unexpected error when decoding result to models.Quote"), err)
	}

	return quotes, nil
}
