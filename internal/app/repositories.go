package app

import (
	"github.com/Sanchir01/currency-wallet/internal/feature/user"
	"github.com/Sanchir01/currency-wallet/internal/feature/wallet"
	"github.com/Sanchir01/currency-wallet/pkg/db"
)

type Repository struct {
	UserRepository   *user.Repository
	WalletRepository *wallet.Repository
}

func NewRepository(databases *db.Database) *Repository {
	return &Repository{
		UserRepository:   user.NewRepository(databases.PrimaryDB),
		WalletRepository: wallet.NewRepository(databases.PrimaryDB),
	}
}
