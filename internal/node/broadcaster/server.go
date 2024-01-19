package broadcaster

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/naturalselectionlabs/rss3-node/config"
	"github.com/robfig/cron/v3"
)

type Broadcaster struct {
	config     *config.File
	cron       *cron.Cron
	httpClient *http.Client
}

func (b *Broadcaster) Run(ctx context.Context) error {
	if err := b.Register(ctx); err != nil {
		return fmt.Errorf("register: %w", err)
	}

	_, err := b.cron.AddFunc("@every 1m", func() {
		if err := b.Heartbeat(ctx); err != nil {
			return
		}
	})
	if err != nil {
		return fmt.Errorf("add heartbeat cron job: %w", err)
	}

	b.cron.Start()

	stopchan := make(chan os.Signal, 1)

	signal.Notify(stopchan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-stopchan

	return nil
}

func NewBroadcaster(_ context.Context, config *config.File) (*Broadcaster, error) {
	return &Broadcaster{
		config:     config,
		cron:       cron.New(),
		httpClient: http.DefaultClient,
	}, nil
}