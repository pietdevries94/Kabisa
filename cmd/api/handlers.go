package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) GetQuoteHandler(w http.ResponseWriter, _ *http.Request) {
	quote, err := app.quoteService.GetRandomQuote()
	if err != nil {
		// TODO: log the error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(quote)
	if err != nil {
		// TODO: log the error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
