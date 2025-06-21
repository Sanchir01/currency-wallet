package app

import (
	"github.com/Sanchir01/currency-wallet/internal/feature/events"
	"github.com/Sanchir01/currency-wallet/internal/feature/user"
	"github.com/Sanchir01/currency-wallet/internal/feature/wallet"
	"github.com/Sanchir01/currency-wallet/pkg/db"
	"log/slog"
)

type Repository struct {
	UserRepository   *user.Repository
	WalletRepository *wallet.Repository
	EventRepository  *events.Repository
}

func NewRepository(databases *db.Database, l *slog.Logger) *Repository {
	return &Repository{
		UserRepository:   user.NewRepository(databases.PrimaryDB),
		WalletRepository: wallet.NewRepository(databases.PrimaryDB, l),
		EventRepository:  events.NewRepository(databases.PrimaryDB),
	}
}
