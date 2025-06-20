package app

import (
	"context"
	"fmt"
	"github.com/Sanchir01/currency-wallet/internal/config"
	"github.com/Sanchir01/currency-wallet/pkg/db"
	"github.com/Sanchir01/currency-wallet/pkg/logger"
	grpcapp "github.com/Sanchir01/currency-wallet/pkg/server/grpc"
	walletsv1 "github.com/Sanchir01/wallets-proto/gen/go/wallets"
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
	echangergrpcurl := fmt.Sprintf("%s:%s", cfg.GrpcClients.GRPCExchanger.Host, cfg.GrpcClients.GRPCExchanger.Port)
	exchanger, err := grpcapp.NewClientGRPC(
		l,
		echangergrpcurl,
		cfg.GrpcClients.GRPCExchanger.Timeout,
		cfg.GrpcClients.GRPCExchanger.Retries,
		walletsv1.NewExchangeServiceClient,
	)

	repo := NewRepository(database, l)
	srv := NewServices(repo, database, l, exchanger)
	handlers := NewHandlers(srv, l)

	return &App{
		Cfg:      cfg,
		Lg:       l,
		DB:       database,
		Handlers: handlers,
	}, nil
}
