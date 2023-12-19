package worker

import (
	"fmt"

	"github.com/naturalselectionlabs/rss3-node/internal/database"
	"github.com/naturalselectionlabs/rss3-node/internal/engine"
	"github.com/naturalselectionlabs/rss3-node/internal/engine/worker/contract/mirror"
	"github.com/naturalselectionlabs/rss3-node/internal/engine/worker/fallback"
	"github.com/naturalselectionlabs/rss3-node/internal/engine/worker/farcaster"
)

func New(config *engine.Config, databaseClient database.Client) (engine.Worker, error) {
	switch config.Worker {
	case engine.Fallback:
		return fallback.NewWorker(config)
	case engine.Mirror:
		return mirror.NewWorker(config, databaseClient)
	case engine.Farcaster:
		return farcaster.NewWorker()
	default:
		return nil, fmt.Errorf("unsupported worker %s", config.Worker)
	}
}
