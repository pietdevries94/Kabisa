package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pietdevries94/Kabisa/api"
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

func (app *application) GetRandomQuote(_ context.Context) (api.GetRandomQuoteRes, error) {
	quote, err := app.quoteService.GetRandomQuote()
	if err != nil {
		app.logger.Error().Err(err).Msg("unexpected error when calling quoteService.GetRandomQuote")
		return app.internalServerError()
	}

	result := &api.Quote{
		ID:     float64(quote.ID),
		Quote:  quote.Quote,
		Author: quote.Author,
	}
	return result, nil
}

func (app *application) internalServerError() (*api.InternalServerErrror, error) {
	return &api.InternalServerErrror{
		Message: "unknown error",
	}, nil
}
