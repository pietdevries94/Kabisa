package models

import "github.com/google/uuid"

type QuoteGame struct {
	ID      uuid.UUID
	Quotes  []*QuoteWithoutAuthor
	Authors []string
}
