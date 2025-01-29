package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pietdevries94/Kabisa/repositories"
	"github.com/pietdevries94/Kabisa/services"
)

type application struct {
	quoteService services.QuoteService
}

func main() {
	app := newApplication()

	r := chi.NewRouter()
	r.Get("/quote", app.GetQuoteHandler)

	// TODO: Make this configurable and log it
	err := http.ListenAndServe("127.0.0.1:3333", r)
	if err != nil {
		// TODO: use the actual logger
		log.Fatal(err)
	}
}

func newApplication() *application {
	client := &http.Client{
		// TODO: Make this configurable
		Timeout: 10 * time.Second,
	}
	dummyJsonRepo := repositories.NewDummyJsonRepo(client)
	quoteService := services.NewQuoteService(dummyJsonRepo)

	return &application{
		quoteService: quoteService,
	}
}
