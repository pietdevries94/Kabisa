package repositories

import (
	"strings"

	"github.com/google/uuid"
	"github.com/pietdevries94/Kabisa/models"
	"github.com/rs/zerolog"
	"golang.org/x/exp/slices"
)

type QuoteGameRepo struct {
	logger *zerolog.Logger
}

// NewQuoteGameRepo returns a new QuoteGameRepo, which creates a manages instances of the quote game.
func NewQuoteGameRepo(logger *zerolog.Logger) *QuoteGameRepo {
	return &QuoteGameRepo{
		logger: logger,
	}
}

// CreateQuoteGame builds a new QuoteGame struct from the given quotes, stores the game in the database for later retrieval and returns the struct
// To make a QuoteGame, the function splits the quotes from the authors and sorts them both alphabetically. As id, it uses an uuid, so players can't
// influence each other's games by guessing valid ids.
func (repo *QuoteGameRepo) CreateQuoteGame(quotes []*models.Quote) (*models.QuoteGame, error) {
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

	// TODO: store the game in the db

	return game, nil
}
