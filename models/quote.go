package models

type Quote struct {
	ID     int
	Quote  string
	Author string
}

type QuoteWithoutAuthor struct {
	ID    int
	Quote string
}
