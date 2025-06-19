package wallet

import (
	"context"
	"github.com/Sanchir01/currency-wallet/internal/domain/models"
	"github.com/Sanchir01/currency-wallet/pkg/api"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type CurrencyWalletResponse struct {
	Rates map[string]float32 `json:"rates"`
}
type HandlerWallets interface {
	GetCurrencyWallets(ctx context.Context) (*models.CurrencyWallet, error)
}

type Handler struct {
	s   HandlerWallets
	log *slog.Logger
}

func NewHandler(s HandlerWallets, log *slog.Logger) *Handler {
	return &Handler{
		s:   s,
		log: log,
	}
}

// @Summary Auth
// @Tags wallet
// @Description register user
// @Accept json
// @Produce json
// @Success 200 {object}  CurrencyWalletResponse
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Security refreshToken
// @Router /exchanger [get]
func (h *Handler) GetAllCurrency(w http.ResponseWriter, r *http.Request) {
	const op = "Wallet.Handler.GetAllCurrency"
	log := h.log.With(slog.String("op", op))
	data, err := h.s.GetCurrencyWallets(r.Context())
	if err != nil {
		log.Error("failed get ", err.Error())
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed get currency wallet"))
		return
	}
	log.Info("Successfully fetched currency data")
	render.JSON(w, r, &CurrencyWalletResponse{
		Rates: data.Balances,
	})
}
