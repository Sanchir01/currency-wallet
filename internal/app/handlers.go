package app

import (
	"github.com/Sanchir01/currency-wallet/internal/feature/user"
	"github.com/Sanchir01/currency-wallet/internal/feature/wallet"
	"log/slog"
)

type Handlers struct {
	UserHandler   *user.Handler
	WalletHandler *wallet.Handler
}

func NewHandlers(services *Services, log *slog.Logger) *Handlers {
	return &Handlers{
		UserHandler:   user.NewHandler(services.UserService, log),
		WalletHandler: wallet.NewHandler(services.WalletService, log),
	}
}
