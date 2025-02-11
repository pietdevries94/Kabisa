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

// A basic quote.
// Ref: #/components/schemas/Quote
type Quote struct {
	ID     int    `json:"id"`
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

// GetID returns the value of ID.
func (s *Quote) GetID() int {
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
func (s *Quote) SetID(val int) {
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

// An answer to the quote game.
// Ref: #/components/schemas/QuoteGameAnswer
type QuoteGameAnswer struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
}

// GetID returns the value of ID.
func (s *QuoteGameAnswer) GetID() int {
	return s.ID
}

// GetAuthor returns the value of Author.
func (s *QuoteGameAnswer) GetAuthor() string {
	return s.Author
}

// SetID sets the value of ID.
func (s *QuoteGameAnswer) SetID(val int) {
	s.ID = val
}

// SetAuthor sets the value of Author.
func (s *QuoteGameAnswer) SetAuthor(val string) {
	s.Author = val
}

// The result of a quote game.
// Ref: #/components/schemas/QuoteGameResult
type QuoteGameResult struct {
	ID      UUID                         `json:"id"`
	Answers []QuoteGameResultAnswersItem `json:"answers"`
}

// GetID returns the value of ID.
func (s *QuoteGameResult) GetID() UUID {
	return s.ID
}

// GetAnswers returns the value of Answers.
func (s *QuoteGameResult) GetAnswers() []QuoteGameResultAnswersItem {
	return s.Answers
}

// SetID sets the value of ID.
func (s *QuoteGameResult) SetID(val UUID) {
	s.ID = val
}

// SetAnswers sets the value of Answers.
func (s *QuoteGameResult) SetAnswers(val []QuoteGameResultAnswersItem) {
	s.Answers = val
}

func (*QuoteGameResult) submitAnswerForQuoteGameRes() {}

type QuoteGameResultAnswersItem struct {
	ID           int    `json:"id"`
	Correct      bool   `json:"correct"`
	ActualAuthor string `json:"actual_author"`
}

// GetID returns the value of ID.
func (s *QuoteGameResultAnswersItem) GetID() int {
	return s.ID
}

// GetCorrect returns the value of Correct.
func (s *QuoteGameResultAnswersItem) GetCorrect() bool {
	return s.Correct
}

// GetActualAuthor returns the value of ActualAuthor.
func (s *QuoteGameResultAnswersItem) GetActualAuthor() string {
	return s.ActualAuthor
}

// SetID sets the value of ID.
func (s *QuoteGameResultAnswersItem) SetID(val int) {
	s.ID = val
}

// SetCorrect sets the value of Correct.
func (s *QuoteGameResultAnswersItem) SetCorrect(val bool) {
	s.Correct = val
}

// SetActualAuthor sets the value of ActualAuthor.
func (s *QuoteGameResultAnswersItem) SetActualAuthor(val string) {
	s.ActualAuthor = val
}

// QuoteWithoutAuthor is used by the quote game.
// Ref: #/components/schemas/QuoteWithoutAuthor
type QuoteWithoutAuthor struct {
	ID    int    `json:"id"`
	Quote string `json:"quote"`
}

// GetID returns the value of ID.
func (s *QuoteWithoutAuthor) GetID() int {
	return s.ID
}

// GetQuote returns the value of Quote.
func (s *QuoteWithoutAuthor) GetQuote() string {
	return s.Quote
}

// SetID sets the value of ID.
func (s *QuoteWithoutAuthor) SetID(val int) {
	s.ID = val
}

// SetQuote sets the value of Quote.
func (s *QuoteWithoutAuthor) SetQuote(val string) {
	s.Quote = val
}

type R404 struct {
	Message string `json:"message"`
}

// GetMessage returns the value of Message.
func (s *R404) GetMessage() string {
	return s.Message
}

// SetMessage sets the value of Message.
func (s *R404) SetMessage(val string) {
	s.Message = val
}

func (*R404) submitAnswerForQuoteGameRes() {}

type R422 struct {
	Errors  []R422ErrorsItem `json:"errors"`
	Message string           `json:"message"`
}

// GetErrors returns the value of Errors.
func (s *R422) GetErrors() []R422ErrorsItem {
	return s.Errors
}

// GetMessage returns the value of Message.
func (s *R422) GetMessage() string {
	return s.Message
}

// SetErrors sets the value of Errors.
func (s *R422) SetErrors(val []R422ErrorsItem) {
	s.Errors = val
}

// SetMessage sets the value of Message.
func (s *R422) SetMessage(val string) {
	s.Message = val
}

func (*R422) submitAnswerForQuoteGameRes() {}

type R422ErrorsItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// GetField returns the value of Field.
func (s *R422ErrorsItem) GetField() string {
	return s.Field
}

// GetMessage returns the value of Message.
func (s *R422ErrorsItem) GetMessage() string {
	return s.Message
}

// SetField sets the value of Field.
func (s *R422ErrorsItem) SetField(val string) {
	s.Field = val
}

// SetMessage sets the value of Message.
func (s *R422ErrorsItem) SetMessage(val string) {
	s.Message = val
}

type R500 struct {
	Message string `json:"message"`
}

// GetMessage returns the value of Message.
func (s *R500) GetMessage() string {
	return s.Message
}

// SetMessage sets the value of Message.
func (s *R500) SetMessage(val string) {
	s.Message = val
}

func (*R500) createNewQuoteGameRes()       {}
func (*R500) getRandomQuoteRes()           {}
func (*R500) submitAnswerForQuoteGameRes() {}

type UUID string
