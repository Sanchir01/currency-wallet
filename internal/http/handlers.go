package httphandlers

import (
	_ "github.com/Sanchir01/currency-wallet/docs"
	"github.com/Sanchir01/currency-wallet/internal/app"
	"github.com/Sanchir01/currency-wallet/internal/http/customiddleware"
	"github.com/Sanchir01/currency-wallet/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log/slog"
	"net/http"
)

func StartHTTTPHandlers(handlers *app.Handlers, domain string, l *slog.Logger) http.Handler {
	router := chi.NewRouter()
	custommiddleware(router, l)
	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", handlers.UserHandler.RegisterHandler)
		r.Post("/login", handlers.UserHandler.LoginHandler)
		r.Group(func(r chi.Router) {
			r.Use(customiddleware.AuthMiddleware(domain))
			r.Get("/exchange/rates", handlers.WalletHandler.GetAllCurrencyHandler)
			r.Get("/balance", handlers.WalletHandler.GetBalanceHandler)
			r.Post("/deposit", handlers.WalletHandler.DepositWallet)
			r.Post("/withdraw", handlers.WalletHandler.WithdrawWallet)
		})
	})
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	return router
}
func custommiddleware(router *chi.Mux, l *slog.Logger) {
	router.Use(middleware.RequestID, middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(logger.NewMiddlewareLogger(l))
	router.Use(customiddleware.PrometheusMiddleware)
}
func StartPrometheusHandlers() http.Handler {
	router := chi.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	return router
}
