package monitor

import (
	"context"
	"fmt"
	"sync"

	"github.com/rss3-network/node/config"
	"go.uber.org/zap"
)

func (m *Monitor) MonitorMockWorkerStatus(ctx context.Context, currentState CheckpointState, latestState uint64) error {
	var wg sync.WaitGroup

	errChan := make(chan error, len(m.config.Node.Decentralized))

	for _, w := range m.config.Node.Decentralized {
		wg.Add(1)

		go func(w *config.Module) {
			defer wg.Done()

			if err := m.processMockWorker(ctx, w, currentState, latestState); err != nil {
				errChan <- err
			}
		}(w)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// processWorker processes the worker status.
func (m *Monitor) processMockWorker(ctx context.Context, w *config.Module, currentState CheckpointState, latestWorkerState uint64) error {
	// get current indexing block height, number or event id and the latest block height, number, timestamp of network
	currentWorkerState, _, err := m.getWorkerIndexingStateByClients(ctx, w.Network, currentState)
	if err != nil {
		zap.L().Error("get latest block height or number", zap.Error(err))
		return err
	}

	// check worker's current status, and flag it as unhealthy if it's left behind the latest block height/number by more than the tolerance
	if err := m.flagUnhealthyWorker(ctx, w.Network.String(), w.ID(), w.Worker.String(), currentWorkerState, latestWorkerState, NetworkTorlerance[w.Network]); err != nil {
		return fmt.Errorf("detect unhealthy: %w", err)
	}

	//	update worker progress (state)
	if err := m.UpdateWorkerProgress(ctx, w.ID(), currentWorkerState); err != nil {
		return fmt.Errorf("update worker progress: %w", err)
	}

	return nil
}
