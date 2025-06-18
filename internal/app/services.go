package app

import (
	"github.com/Sanchir01/currency-wallet/internal/feature/user"
	"github.com/Sanchir01/currency-wallet/pkg/db"
	"log/slog"
)

type Services struct {
	UserService *user.Service
}

func NewServices(repos *Repository, db *db.Database, l *slog.Logger) *Services {
	return &Services{
		UserService: user.NewService(repos.UserRepository, db.PrimaryDB, l),
	}
}
