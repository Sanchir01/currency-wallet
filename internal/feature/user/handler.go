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

type Handler struct {
	Service HandlerUser
	log     *slog.Logger
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=HandlerUser
type HandlerUser interface {
	Register(ctx context.Context, email, username, password string) (*uuid.UUID, error)
	Login(ctx context.Context, email, password string) (*DatabaseUser, error)
}

func NewHandler(s HandlerUser, lg *slog.Logger) *Handler {
	return &Handler{
		Service: s,
		log:     lg,
	}
}

// @Summary Auth
// @Tags auth
// @Description register user
// @Accept json
// @Produce json
// @Param input body LoginRequest true "login body"
// @Success 201 {object}  AuthResponse
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Router /register [post]
func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.register"
	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	var req AuthRequest
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
	id, err := h.Service.Register(r.Context(), req.Email, req.Username, req.Password)
	if errors.Is(err, utils.ErrorUserAlreadyExists) {
		log.Error("user already exists", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("Пользователь с таким email или username уже существует"))
		return
	}
	if err != nil {
		log.Error("failed to register user", logger.Err(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("ошибка регистрации"))
		return
	}
	log.Info("login success")

	if err = AddCookieTokens(*id, w, "localhost"); err != nil {
		log.Error("register cookie errors", err.Error())
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed to register cookie"))
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, AuthResponse{
		Response: api.OK(),
	})
}

// @Summary Login
// @Tags auth
// @Description login user
// @Accept json
// @Produce json
// @Param input body LoginRequest true "auth body"
// @Success 200 {object}  LoginResponse
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Router /login [post]
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	const op = "User.Handler.Login"
	log := h.log.With(slog.String("op", op))
	var req LoginRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", slog.Any("err", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("Ошибка при валидации данных"))
		return
	}
	if err := validator.New().Struct(req); err != nil {
		log.Error("invalid request", logger.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	user, err := h.Service.Login(r.Context(), req.Email, req.Password)
	if errors.Is(err, utils.ErrorUserNotFound) {
		log.Error("invalid login", logger.Err(err))
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, api.Error("Неправильный логин"))
		return
	}
	if errors.Is(err, utils.ErrorInvalidPassword) {
		log.Error("invalid password", logger.Err(err))
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, api.Error("Неправильный логин"))
		return
	}
	if err != nil {
		log.Error("failed to login", logger.Err(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("internal server error"))
		return
	}
	render.JSON(w, r, LoginResponse{
		Response: api.OK(),
		Email:    user.Email,
		Username: user.Name,
	})
}
