package user

import (
	"context"
	"errors"

	"github.com/Sanchir01/currency-wallet/pkg/logger"
	"github.com/Sanchir01/currency-wallet/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log/slog"
	"net/http"

	"github.com/Sanchir01/currency-wallet/pkg/api"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type AuthRequest struct {
	Email    string `json:"email" validate:"required"`
	Username string `json:"username" validate:"required,min=1,max=100"`
	Password string `json:"password" validate:"required,min=6"`
}
type SendCoinsRequest struct {
	Email string `json:"toUser" validate:"required,email"`
	Coins int64  `json:"amount" validate:"required"`
}
type BuyProductResponse struct {
	api.Response
	Ok string `json:"ok"`
}
type AuthResponse struct {
	api.Response
}
type AllUserCoinsInfoResponse struct {
	api.Response
	GetAllUserCoinsInfo *GetAllUserCoinsInfo
}
type Handler struct {
	Service HandlerUser
	Log     *slog.Logger
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=HandlerUser
type HandlerUser interface {
	Register(ctx context.Context, email, username, password string) (*uuid.UUID, error)
}

func NewHandler(s HandlerUser, lg *slog.Logger) *Handler {
	return &Handler{
		Service: s,
		Log:     lg,
	}
}

// @Summary Auth
// @Tags user
// @Description buy product endpoin
// @Accept json
// @Produce json
// @Param input body AuthRequest true "auth body"
// @Success 200 {object}  AuthResponse
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Router /api/auth [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.register"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	var req AuthRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", slog.Any("err", err))
		render.JSON(w, r, api.Error("Ошибка при валидации данных"))
		return
	}
	if err := validator.New().Struct(req); err != nil {
		log.Error("invalid request", logger.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	id, err := h.Service.Register(r.Context(), req.Email, req.Username, req.Password)
	if errors.Is(err, utils.ErrorUserAlreadyExists) {
		log.Error("user already exists", logger.Err(err))
		render.JSON(w, r, api.Error("Пользователь с таким email или username уже существует"))
		return
	}
	if err != nil {
		log.Error("failed to register user", logger.Err(err))
		render.JSON(w, r, api.Error("ошибка регистрации"))
		return
	}
	log.Info("login success")

	if err = AddCookieTokens(*id, w, "localhost"); err != nil {
		log.Error("register cookie errors", err.Error())
		render.JSON(w, r, api.Error("failed to register cookie"))
	}
	render.JSON(w, r, AuthResponse{
		Response: api.OK(),
	})
}
