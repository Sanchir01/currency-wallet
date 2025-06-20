package wallet

import (
	"context"
	contextkey "github.com/Sanchir01/currency-wallet/internal/domain/contants"
	"github.com/Sanchir01/currency-wallet/internal/domain/models"
	httphandlers "github.com/Sanchir01/currency-wallet/internal/http/customiddleware"
	"github.com/Sanchir01/currency-wallet/pkg/api"
	"github.com/Sanchir01/currency-wallet/pkg/logger"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type CurrencyWalletResponse struct {
	Rates map[string]float32 `json:"rates"`
}
type HandlerWallets interface {
	GetCurrencyWallets(ctx context.Context) (*models.CurrencyWallet, error)
	GetBalance(ctx context.Context, id uuid.UUID) (*models.CurrencyWallet, error)
	WalletDepositOrWithDraw(ctx context.Context, id uuid.UUID, currency string, amount float32, typedepo contextkey.OperationType) (*models.CurrencyWallet, error)
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

// @Summary GetAllCurrencyHandler
// @Tags wallet
// @Description all currency wallet
// @Accept json
// @Produce json
// @Success 200 {object}  CurrencyWalletResponse
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Security refreshToken
// @Router /exchanger [get]
func (h *Handler) GetAllCurrencyHandler(w http.ResponseWriter, r *http.Request) {
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

// @Summary GetBalanceHandler
// @Tags wallet
// @Description balance user
// @Accept json
// @Produce json
// @Success 200 {object}  CurrencyWalletResponse
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Security refreshToken
// @Router /balance [get]
func (h *Handler) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	const op = "Wallet.Handler.GetAllCurrency"
	log := h.log.With(slog.String("op", op))
	userid, err := httphandlers.GetJWTClaimsFromCtx(r.Context())
	if err != nil {
		log.Error("failed get user id from jwt", err.Error())
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, api.Error("Unauthorized"))
		return
	}
	data, err := h.s.GetBalance(r.Context(), userid.ID)
	if err != nil {
		log.Error("failed get currency balance", err.Error())
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed get currency balance"))
		return
	}
	render.JSON(w, r, &CurrencyWalletResponse{
		Rates: data.Balances,
	})
}

// @Summary GetBalanceHandler
// @Tags wallet
// @Description balance user
// @Accept json
// @Produce json
// @Param input body DepositOrWithdrawRequest true "deposit body"
// @Success 200 {object}  CurrencyWalletResponse
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Security refreshToken
// @Router /deposit [post]
func (h *Handler) DepositWallet(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.DepositWallet"
	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	userid, err := httphandlers.GetJWTClaimsFromCtx(r.Context())
	if err != nil {
		log.Error("failed get user id from jwt", err.Error())
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, api.Error("Unauthorized"))
		return
	}
	var req DepositOrWithdrawRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", slog.Any("err", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("Ошибка при валидации данных"))
		return
	}
	if err := validator.New().Struct(req); err != nil {
		log.Error("invalid request", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	data, err := h.s.WalletDepositOrWithDraw(r.Context(), userid.ID, req.Currency, req.Amount, contextkey.OperationTypeDeposit)
	if err != nil {
		log.Error("failed deposit currency wallet", err.Error())
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed deposit currency wallet"))
		return
	}
	render.JSON(w, r, data)
}

// @Summary WithdrawWallet
// @Tags wallet
// @Description balance user
// @Accept json
// @Produce json
// @Param input body DepositOrWithdrawRequest true "deposit body"
// @Success 200 {object}  CurrencyWalletResponse
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Security refreshToken
// @Router /deposit [post]
func (h *Handler) WithdrawWallet(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.WithdrawWallet"
	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	userid, err := httphandlers.GetJWTClaimsFromCtx(r.Context())
	if err != nil {
		log.Error("failed get user id from jwt", err.Error())
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, api.Error("Unauthorized"))
		return
	}
	var req DepositOrWithdrawRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", slog.Any("err", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("Ошибка при валидации данных"))
		return
	}
	if err := validator.New().Struct(req); err != nil {
		log.Error("invalid request", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	data, err := h.s.WalletDepositOrWithDraw(r.Context(), userid.ID, req.Currency, req.Amount, contextkey.OperationTypeWithdraw)
	if err != nil {
		log.Error("failed deposit currency wallet", err.Error())
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed deposit currency wallet"))
		return
	}
	render.JSON(w, r, data)
}
