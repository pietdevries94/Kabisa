package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/pietdevries94/Kabisa/models"
	"github.com/pietdevries94/Kabisa/openapi"
)

func (app *application) GetRandomQuote(ctx context.Context) (openapi.GetRandomQuoteRes, error) {
	quote, err := app.quoteService.GetRandomQuote(ctx)
	if err != nil {
		app.logger.Error().Err(err).Msg("unexpected error when calling quoteService.GetRandomQuote")
		return app.internalServerError()
	}

	result := &openapi.Quote{
		ID:     quote.ID,
		Quote:  quote.Quote,
		Author: quote.Author,
	}
	return result, nil
}

func (app *application) CreateNewQuoteGame(ctx context.Context) (openapi.CreateNewQuoteGameRes, error) {
	game, err := app.quoteService.CreateQuoteGame(ctx)
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
			ID:    q.ID,
			Quote: q.Quote,
		}
	}

	return result, nil
}

func (app *application) SubmitAnswerForQuoteGame(ctx context.Context, answers []openapi.QuoteGameAnswer, params openapi.SubmitAnswerForQuoteGameParams) (openapi.SubmitAnswerForQuoteGameRes, error) {
	id, err := uuid.Parse(string(params.ID))
	if err != nil {
		return app.notFound()
	}

	answerMap := make(models.QuoteGameAnswerMap)
	for _, a := range answers {
		answerMap[a.ID] = a.Author
	}

	gameResult, err := app.quoteService.SubmitAnswerToQuoteGame(ctx, id, answerMap)
	if err != nil {
		return app.unprocessableContent(err)
	}

	result := &openapi.QuoteGameResult{
		ID:      openapi.UUID(gameResult.ID.String()),
		Answers: make([]openapi.QuoteGameResultAnswersItem, len(gameResult.Answers)),
	}
	for i, a := range gameResult.Answers {
		result.Answers[i] = openapi.QuoteGameResultAnswersItem{
			ID:           a.ID,
			Correct:      a.Correct,
			ActualAuthor: a.Author,
		}
	}

	return result, nil
}

func (app *application) unprocessableContent(err error) (openapi.SubmitAnswerForQuoteGameRes, error) {
	pe, ok := err.(*models.PublicError)
	if !ok {
		return app.internalServerError()
	}
	return &openapi.R422{
		Message: pe.Error(),
	}, nil
}

func (app *application) internalServerError() (*openapi.R500, error) {
	return &openapi.R500{
		Message: "unknown error",
	}, nil
}

func (app *application) notFound() (*openapi.R404, error) {
	return &openapi.R404{
		Message: "not_found",
	}, nil
}
