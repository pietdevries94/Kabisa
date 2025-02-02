// Code generated by ogen, DO NOT EDIT.

package openapi

type CreateNewQuoteGameOK struct {
	ID      UUID                 `json:"id"`
	Quotes  []QuoteWithoutAuthor `json:"quotes"`
	Authors []string             `json:"authors"`
}

// GetID returns the value of ID.
func (s *CreateNewQuoteGameOK) GetID() UUID {
	return s.ID
}

// GetQuotes returns the value of Quotes.
func (s *CreateNewQuoteGameOK) GetQuotes() []QuoteWithoutAuthor {
	return s.Quotes
}

// GetAuthors returns the value of Authors.
func (s *CreateNewQuoteGameOK) GetAuthors() []string {
	return s.Authors
}

// SetID sets the value of ID.
func (s *CreateNewQuoteGameOK) SetID(val UUID) {
	s.ID = val
}

// SetQuotes sets the value of Quotes.
func (s *CreateNewQuoteGameOK) SetQuotes(val []QuoteWithoutAuthor) {
	s.Quotes = val
}

// SetAuthors sets the value of Authors.
func (s *CreateNewQuoteGameOK) SetAuthors(val []string) {
	s.Authors = val
}

func (*CreateNewQuoteGameOK) createNewQuoteGameRes() {}

type InternalServerErrror struct {
	Message string `json:"message"`
}

// GetMessage returns the value of Message.
func (s *InternalServerErrror) GetMessage() string {
	return s.Message
}

// SetMessage sets the value of Message.
func (s *InternalServerErrror) SetMessage(val string) {
	s.Message = val
}

func (*InternalServerErrror) createNewQuoteGameRes() {}
func (*InternalServerErrror) getRandomQuoteRes()     {}

// A basic quote.
// Ref: #/components/schemas/Quote
type Quote struct {
	ID     float64 `json:"id"`
	Quote  string  `json:"quote"`
	Author string  `json:"author"`
}

// GetID returns the value of ID.
func (s *Quote) GetID() float64 {
	return s.ID
}

// GetQuote returns the value of Quote.
func (s *Quote) GetQuote() string {
	return s.Quote
}

// GetAuthor returns the value of Author.
func (s *Quote) GetAuthor() string {
	return s.Author
}

// SetID sets the value of ID.
func (s *Quote) SetID(val float64) {
	s.ID = val
}

// SetQuote sets the value of Quote.
func (s *Quote) SetQuote(val string) {
	s.Quote = val
}

// SetAuthor sets the value of Author.
func (s *Quote) SetAuthor(val string) {
	s.Author = val
}

func (*Quote) getRandomQuoteRes() {}

// QuoteWithoutAuthor is used by the quote game.
// Ref: #/components/schemas/QuoteWithoutAuthor
type QuoteWithoutAuthor struct {
	ID    float64 `json:"id"`
	Quote string  `json:"quote"`
}

// GetID returns the value of ID.
func (s *QuoteWithoutAuthor) GetID() float64 {
	return s.ID
}

// GetQuote returns the value of Quote.
func (s *QuoteWithoutAuthor) GetQuote() string {
	return s.Quote
}

// SetID sets the value of ID.
func (s *QuoteWithoutAuthor) SetID(val float64) {
	s.ID = val
}

// SetQuote sets the value of Quote.
func (s *QuoteWithoutAuthor) SetQuote(val string) {
	s.Quote = val
}

type UUID string
