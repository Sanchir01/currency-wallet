package app

import (
	"context"
	"github.com/Sanchir01/currency-wallet/internal/config"
	"github.com/Sanchir01/currency-wallet/pkg/db"
	"github.com/Sanchir01/currency-wallet/pkg/logger"
	"log/slog"
)

type App struct {
	Cfg      *config.Config
	Lg       *slog.Logger
	DB       *db.Database
	Handlers *Handlers
}

func NewApp(ctx context.Context) (*App, error) {
	cfg := config.InitConfig()
	l := logger.SetupLogger(cfg.Env)
	database, err := db.NewDataBases(cfg, ctx)
	if err != nil {
		return nil, err
	}
	repo := NewRepository(database)
	srv := NewServices(repo, database, l)
	handlers := NewHandlers(srv, l)
	return &App{
		Cfg:      cfg,
		Lg:       l,
		DB:       database,
		Handlers: handlers,
	}, nil
}
