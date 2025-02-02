package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pietdevries94/Kabisa/openapi"
)

func (app *application) GetQuoteHandler(w http.ResponseWriter, _ *http.Request) {
	quote, err := app.quoteService.GetRandomQuote()
	if err != nil {
		app.logger.Error().Err(err).Msg("unexpected error when calling quoteService.GetRandomQuote")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(quote)
	if err != nil {
		app.logger.Error().Err(err).Msg("unexpected error when encoding quote handler to json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *application) GetRandomQuote(_ context.Context) (openapi.GetRandomQuoteRes, error) {
	quote, err := app.quoteService.GetRandomQuote()
	if err != nil {
		app.logger.Error().Err(err).Msg("unexpected error when calling quoteService.GetRandomQuote")
		return app.internalServerError()
	}

	result := &openapi.Quote{
		ID:     float64(quote.ID),
		Quote:  quote.Quote,
		Author: quote.Author,
	}
	return result, nil
}

func (app *application) CreateNewQuoteGame(_ context.Context) (openapi.CreateNewQuoteGameRes, error) {
	game, err := app.quoteService.CreateQuoteGame()
	if err != nil {
		app.logger.Error().Err(err).Msg("unexpected error when calling quoteService.CreateQuoteGame")
		return app.internalServerError()
	}

	result := &openapi.CreateNewQuoteGameOK{
		ID:      openapi.UUID(game.ID.String()),
		Authors: game.Authors,
	}
	result.Quotes = make([]openapi.QuoteWithoutAuthor, len(game.Quotes))
	for i, q := range game.Quotes {
		result.Quotes[i] = openapi.QuoteWithoutAuthor{
			ID:    float64(q.ID),
			Quote: q.Quote,
		}
	}

	return result, nil
}

func (app *application) internalServerError() (*openapi.InternalServerErrror, error) {
	return &openapi.InternalServerErrror{
		Message: "unknown error",
	}, nil
}
