//go:build wireinject
// +build wireinject

package app

import (
	"fmt"
	"log/slog"
	"object-storage/internal/config"
	"object-storage/internal/handlers/grpc"
	"object-storage/internal/storage"
	"os"

	"github.com/google/wire"
	"github.com/nats-io/nats.go"
)

func Init() (*App, func(), error) {
	panic(
		wire.Build(
			newApp,
			wire.NewSet(config.New),
			wire.NewSet(initLogger),
			wire.NewSet(initNats),

			wire.NewSet(storage.New),
			wire.Bind(new(grpc.Storage), new(*storage.ObjectStorage)),

			// handlers
			wire.NewSet(grpc.New),
		),
	)
}

func initLogger(cfg *config.Config) *slog.Logger {
	var level slog.Level

	switch cfg.Logger.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	}

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}

func initNats(cfg *config.Config) (*nats.Conn, func(), error) {
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%d", cfg.Nats.Host, cfg.Nats.Port))
	if err != nil {
		return nil, nil, err
	}
	return nc, func() {
		nc.Close()
	}, nil
}
