package app

import (
	"github.com/Sanchir01/currency-wallet/internal/feature/user"
	"github.com/Sanchir01/currency-wallet/pkg/db"
)

type Repository struct {
	UserRepository *user.Repository
}

func NewRepository(databases *db.Database) *Repository {
	return &Repository{
		UserRepository: user.NewRepository(databases.PrimaryDB),
	}
}
