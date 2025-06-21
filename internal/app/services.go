package app

import (
	"github.com/Sanchir01/currency-wallet/internal/feature/events"
	"github.com/Sanchir01/currency-wallet/internal/feature/user"
	"github.com/Sanchir01/currency-wallet/internal/feature/wallet"
	"github.com/Sanchir01/currency-wallet/pkg/db"
	kafkaclient "github.com/Sanchir01/currency-wallet/pkg/events"
	walletsv1 "github.com/Sanchir01/wallets-proto/gen/go/wallets"
	"log/slog"
)

type Services struct {
	UserService   *user.Service
	WalletService *wallet.Service
	EventService  *events.Service
}

func NewServices(repos *Repository, db *db.Database, l *slog.Logger, exchanger walletsv1.ExchangeServiceClient, producer *kafkaclient.Producer) *Services {
	return &Services{
		UserService:   user.NewService(repos.UserRepository, repos.WalletRepository, db.PrimaryDB, l),
		WalletService: wallet.NewService(repos.WalletRepository, repos.EventRepository, db.PrimaryDB, db.RedisDB, exchanger, l),
		EventService:  events.NewEventService(l, repos.EventRepository, producer),
	}
}
