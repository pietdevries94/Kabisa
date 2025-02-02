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

func (repo *QuoteGameRepo) ValidateIDAndAnswerIDs(ctx context.Context, id uuid.UUID, answers models.QuoteGameAnswerMap) (answerIDs []int, err error) {
	_ = ctx
	_ = id
	_ = answers
	panic("// TODO implement this")
}

func (repo *QuoteGameRepo) ValidateAnswersAndCreateGameResult(ctx context.Context, id uuid.UUID, quotes map[int]*models.Quote, answers models.QuoteGameAnswerMap) (*models.QuoteGameResult, error) {
	_ = ctx
	_ = id
	_ = quotes
	_ = answers
	panic("// TODO implement this")
}
