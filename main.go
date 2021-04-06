package main

import (
	"context"
	"database/sql"

	"github.com/adamgoose/walker-server/lib/walker"
	"github.com/heroiclabs/nakama-common/runtime"
)

// InitModule is the main export expected by the runtime. Register stuff here
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {

	if err := walker.Register(ctx, initializer); err != nil {
		logger.Error("Unable to Register Walker Match Handler: %v", err)
		return err
	}

	return nil
}
