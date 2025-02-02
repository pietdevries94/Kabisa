package models

import "github.com/google/uuid"

type QuoteGame struct {
	ID      uuid.UUID
	Quotes  []*QuoteWithoutAuthor
	Authors []string
}

type QuoteGameAnswerMap map[int]string

type QuoteGameResult struct {
	ID      uuid.UUID
	Answers []*QuoteGameActualAnswer
}

type QuoteGameActualAnswer struct {
	Quote
	Correct bool
}
