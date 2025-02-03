package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pietdevries94/Kabisa/models"
	"github.com/rs/zerolog"
	"golang.org/x/exp/slices"

	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/im"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
	"github.com/stephenafamo/bob/dialect/sqlite/um"
)

type QuoteGameRepo struct {
	logger *zerolog.Logger
	db     *sql.DB
}

// NewQuoteGameRepo returns a new QuoteGameRepo, which creates a manages instances of the quote game.
func NewQuoteGameRepo(logger *zerolog.Logger, db *sql.DB) *QuoteGameRepo {
	return &QuoteGameRepo{
		logger: logger,
		db:     db,
	}
}

// CreateQuoteGame builds a new QuoteGame struct from the given quotes, stores the game in the database for later retrieval and returns the struct
// To make a QuoteGame, the function splits the quotes from the authors and sorts them both alphabetically. As id, it uses an uuid, so players can't
// influence each other's games by guessing valid ids.
func (repo *QuoteGameRepo) CreateQuoteGame(ctx context.Context, quotes []*models.Quote) (*models.QuoteGame, error) {
	// Currently the game only supports three quotes
	if len(quotes) != 3 {
		return nil, fmt.Errorf("number of quotes should be 3. Given: %d", len(quotes))
	}

	// First we prepare a new game
	game := &models.QuoteGame{
		ID:      uuid.New(),
		Quotes:  make([]*models.QuoteWithoutAuthor, len(quotes)),
		Authors: make([]string, len(quotes)),
	}

	// We split the quotes
	for i, q := range quotes {
		game.Quotes[i] = &models.QuoteWithoutAuthor{
			ID:    q.ID,
			Quote: q.Quote,
		}
		game.Authors[i] = q.Author
	}

	// And sort them both
	slices.SortFunc(game.Quotes, func(a *models.QuoteWithoutAuthor, b *models.QuoteWithoutAuthor) int {
		return strings.Compare(a.Quote, b.Quote)
	})
	slices.Sort(game.Authors)

	// Now we build the query to store it in the database
	queryString, args, err := sqlite.Insert(
		im.Into("quote_game", "id", "quote1_id", "quote2_id", "quote3_id", "created_at"),
		im.Values(sqlite.Arg(game.ID, game.Quotes[0].ID, game.Quotes[1].ID, game.Quotes[2].ID, time.Now())),
	).Build(ctx)
	if err != nil {
		repo.logger.Error().Err(err).Msg("could not build query")
		return nil, errors.Join(errors.New("could not build query"), err)
	}

	// Execute the query
	_, err = repo.db.ExecContext(ctx, queryString, args...)
	if err != nil {
		repo.logger.Error().Err(err).Msg("could not execute query")
		return nil, errors.Join(errors.New("could not execute query"), err)
	}

	// Finally, we return the game
	return game, nil
}

// ValidateIDAndAnswerIDs gets the game information from the database, runs a couple checks and returns the quote_ids in order from the database.
// The following checks are performed:
//   - Does the id exist
//   - Is the completed_at null
//   - Is the created_at within a minute ago
//   - Are the quote ids present in the map
//   - Are only the quote ids present in the map
func (repo *QuoteGameRepo) ValidateIDAndAnswerIDs(ctx context.Context, id uuid.UUID, answers models.QuoteGameAnswerMap) (quoteIDs []int, err error) {
	// There have to be exactly three quotes in the map
	if len(answers) != 3 {
		return nil, models.ErrInvalidQuoteID
	}

	queryString, args, err := sqlite.Select(
		sm.From("quote_game"),
		sm.Columns("quote1_id", "quote2_id", "quote3_id", "created_at", "completed_at"),
		sm.Where(sqlite.Quote("id").EQ(sqlite.Arg(id))),
	).Build(ctx)
	if err != nil {
		repo.logger.Error().Err(err).Msg("could not build query")
		return nil, errors.Join(errors.New("could not build query"), err)
	}

	quoteIDs = make([]int, 3)
	var createdAt time.Time
	var completedAt sql.NullTime
	err = repo.db.QueryRowContext(ctx, queryString, args...).Scan(&quoteIDs[0], &quoteIDs[1], &quoteIDs[2], &createdAt, &completedAt)
	if err == sql.ErrNoRows {
		return nil, models.ErrQuoteGameIdNotFound
	}
	if err != nil {
		repo.logger.Error().Err(err).Msg("could not execute query")
		return nil, errors.Join(errors.New("could not execute query"), err)
	}

	// We check if the game is not completed yet
	if completedAt.Valid {
		return nil, models.ErrQuoteGameIdNotFound
	}

	// Or expired
	if time.Now().After(createdAt.Add(time.Minute)) {
		return nil, models.ErrQuoteGameIdNotFound
	}

	// And if the answer ids match with the quote ids
	for _, id := range quoteIDs {
		if _, ok := answers[id]; !ok {
			return nil, models.ErrInvalidQuoteID
		}
	}

	return quoteIDs, nil
}

// ValidateAnswersAndCreateGameResult compares the given answers to the quote authors, compiles a result and puts it in the database.
func (repo *QuoteGameRepo) ValidateAnswersAndCreateGameResult(ctx context.Context, id uuid.UUID, quoteIDs []int, quotes map[int]*models.Quote, answers models.QuoteGameAnswerMap) (*models.QuoteGameResult, error) {
	gameResult := &models.QuoteGameResult{
		ID:      id,
		Answers: make([]*models.QuoteGameActualAnswer, len(quoteIDs)),
	}

	// Compile the actual answers and check if the given answer is correct
	for i, id := range quoteIDs {
		quote := quotes[id]
		correct := quote.Author == answers[id]
		gameResult.Answers[i] = &models.QuoteGameActualAnswer{
			Quote:   *quote,
			Correct: correct,
		}
	}

	// We set the result in the database
	queryString, args, err := sqlite.Update(
		um.Table("quote_game"),
		um.SetCol("quote1_correct").ToArg(gameResult.Answers[0].Correct),
		um.SetCol("quote2_correct").ToArg(gameResult.Answers[1].Correct),
		um.SetCol("quote3_correct").ToArg(gameResult.Answers[2].Correct),
		um.SetCol("completed_at").ToArg(time.Now()),
		um.Where(sqlite.Quote("id").EQ(sqlite.Arg(id))),
	).Build(ctx)
	if err != nil {
		repo.logger.Error().Err(err).Msg("could not build query")
		return nil, errors.Join(errors.New("could not build query"), err)
	}

	// Execute the query
	_, err = repo.db.ExecContext(ctx, queryString, args...)
	if err != nil {
		repo.logger.Error().Err(err).Msg("could not execute query")
		return nil, errors.Join(errors.New("could not execute query"), err)
	}

	// And return the result
	return gameResult, nil
}
