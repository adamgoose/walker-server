package walker

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
)

// ModuleName defines the Nakama module name to use
var ModuleName = "walker"

const (
	OpCodeNoop      = iota // 0
	OpCodeMove             // 1 client -> client
	OpCodeSpawnBeer        // 2 server -> client
	OpCodeClaimBeer        // 3 client -> server, client
	OpCodeScores           // 4 server -> client
)

// Walker represents the Walker Match Handler
type Walker struct{}

// Register registers the Walker Module
func Register(ctx context.Context, init runtime.Initializer) error {
	w := &Walker{}

	if err := init.RegisterMatch(ModuleName, w.registerMatch); err != nil {
		return err
	}

	if err := init.RegisterMatchmakerMatched(w.registerMatchmakerMatched); err != nil {
		return err
	}

	return nil
}

func (w *Walker) registerMatch(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
	return w, nil
}

func (w *Walker) registerMatchmakerMatched(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, entries []runtime.MatchmakerEntry) (string, error) {
	logger.WithField("entries", entries).Info("MatchmakerMatched")
	return nk.MatchCreate(ctx, ModuleName, map[string]interface{}{"users": entries})
}
