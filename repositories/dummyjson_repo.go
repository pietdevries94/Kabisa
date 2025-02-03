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
func (repo *DummyJsonRepo) GetRandomQuotes(ctx context.Context, amount int) ([]*models.Quote, error) {
	if amount < 1 || amount > 10 {
		// The api only accepts 1 to 10
		return nil, fmt.Errorf("amount should be between 1 and 10. Given: %d", amount)
	}

	url := fmt.Sprintf("https://dummyjson.com/quotes/random/%d", amount)

	resp, err := repo.get(ctx, url)
	if err != nil {
		return nil, err
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

// GetQuotes retrieves a map of quotes from an api. The api only supports doing this one by one.
// It does so synchronized, but if this function needs to handle larger numbers, it could be
// refactored to build up the map async with a limit of parallel fetches.
// If any request errors, we return an error and no map.
func (repo *DummyJsonRepo) GetQuotes(ctx context.Context, ids []int) (map[int]*models.Quote, error) {
	m := map[int]*models.Quote{}
	for _, id := range ids {
		quote, err := repo.GetQuote(ctx, id)
		if err != nil {
			return nil, err
		}
		m[id] = quote
	}
	return m, nil
}

// GetQuote gets a quote by id, or returns a public error when not found. Other errors get logged
func (repo *DummyJsonRepo) GetQuote(ctx context.Context, id int) (*models.Quote, error) {
	url := fmt.Sprintf("https://dummyjson.com/quotes/%d", id)

	resp, err := repo.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// we want to return a public error if the response is a 404
	if resp.StatusCode == http.StatusNotFound {
		return nil, models.NewPublicErrorf("unknown_quote_id: %d", id)
	}

	if resp.StatusCode != http.StatusOK {
		repo.logger.Error().Int("status code", resp.StatusCode).Msg("unexpected status code received")
		return nil, fmt.Errorf("unexpected status code received: %d", resp.StatusCode)
	}

	var quote *models.Quote
	err = json.NewDecoder(resp.Body).Decode(&quote)
	if err != nil {
		repo.logger.Error().Err(err).Msg("unexpected error when decoding result to models.Quote")
		return nil, errors.Join(errors.New("unexpected error when decoding result to models.Quote"), err)
	}
	return quote, err
}

func (repo *DummyJsonRepo) get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		repo.logger.Error().Err(err).Msg("unexpected error when creating request for retrieving random quote from api")
		return nil, errors.Join(errors.New("unexpected error when creating request for retrieving random quote from api"), err)
	}

	resp, err := repo.httpClient.Do(req)
	if err != nil {
		repo.logger.Error().Err(err).Msg("unexpected error when retrieving random quote from api")
		return nil, errors.Join(errors.New("unexpected error when retrieving random quote from api"), err)
	}

	if resp.Body == nil {
		repo.logger.Error().Err(err).Int("status code", resp.StatusCode).Msg("no body received")
		return nil, errors.New("no body received")
	}

	return resp, nil
}
