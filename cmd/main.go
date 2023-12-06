package main

import (
	"context"
	"fmt"

	"github.com/naturalselectionlabs/rss3-node/internal/config"
	"github.com/naturalselectionlabs/rss3-node/internal/config/flag"
	"github.com/naturalselectionlabs/rss3-node/internal/constant"
	"github.com/naturalselectionlabs/rss3-node/internal/database"
	"github.com/naturalselectionlabs/rss3-node/internal/database/dialer"
	"github.com/naturalselectionlabs/rss3-node/internal/engine"
	"github.com/naturalselectionlabs/rss3-node/internal/node/indexer"
	"github.com/naturalselectionlabs/rss3-node/schema/filter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var command = cobra.Command{
	Use:           constant.Name,
	Version:       constant.BuildVersion(),
	SilenceUsage:  true,
	SilenceErrors: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		config, err := config.Setup(viper.GetString(flag.KeyConfig))
		if err != nil {
			return fmt.Errorf("setup config file: %w", err)
		}

		// Dial and migrate database.
		databaseClient, err := dialer.Dial(cmd.Context(), config.Database)
		if err != nil {
			return fmt.Errorf("dial database: %w", err)
		}

		if err := databaseClient.Migrate(cmd.Context()); err != nil {
			return fmt.Errorf("migrate database: %w", err)
		}

		switch viper.GetString(flag.KeyModule) {
		case "explorer":
		case "indexer":
			return runIndexer(cmd.Context(), config, databaseClient)
		}

		return fmt.Errorf("unsupported module %s", viper.GetString(flag.KeyModule))
	},
}

func runIndexer(ctx context.Context, config *config.File, databaseClient database.Client) error {
	network, err := filter.NetworkString(viper.GetString(flag.KeyIndexerNetwork))
	if err != nil {
		return fmt.Errorf("network string: %w", err)
	}

	chain, err := filter.ChainString(network, viper.GetString(flag.KeyIndexerChain))
	if err != nil {
		return fmt.Errorf("chain string: %w", err)
	}

	worker, err := engine.NameString(viper.GetString(flag.KeyIndexerWorker))
	if err != nil {
		return fmt.Errorf("worker string: %w", err)
	}

	for _, nodeConfig := range config.Node.Decentralized {
		if nodeConfig.Chain == chain && nodeConfig.Worker == worker {
			server, err := indexer.NewServer(ctx, nodeConfig, databaseClient)
			if err != nil {
				return fmt.Errorf("new server: %w", err)
			}

			return server.Run(ctx)
		}
	}

	return fmt.Errorf("unsupported indexer %s.%s.%s", network, chain, worker)
}

func initializeLogger() {
	if viper.GetString(config.Environment) == config.EnvironmentDevelopment {
		zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	} else {
		zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	}
}

func init() {
	initializeLogger()

	command.PersistentFlags().String(flag.KeyConfig, "./deploy/config.development.yaml", "config file path")
	command.PersistentFlags().String(flag.KeyModule, "indexer", "module name")
	command.PersistentFlags().String(flag.KeyIndexerNetwork, "ethereum", "indexer network")
	command.PersistentFlags().String(flag.KeyIndexerChain, "mainnet", "indexer chain")
	command.PersistentFlags().String(flag.KeyIndexerWorker, "fallback", "indexer worker")
}

func main() {
	if err := command.ExecuteContext(context.Background()); err != nil {
		zap.L().Fatal("execute command", zap.Error(err))
	}
}
