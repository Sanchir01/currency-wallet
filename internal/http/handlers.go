package httphandlers

import (
	_ "github.com/Sanchir01/currency-wallet/docs"
	"github.com/Sanchir01/currency-wallet/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
)

func StartHTTTPHandlers(handlers *app.Handlers, domain string) http.Handler {
	router := chi.NewRouter()
	custommiddleware(router, domain)
	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/register", handlers.UserHandler.Register)
	})
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	return router
}
func custommiddleware(router *chi.Mux, domain string) {
	router.Use(middleware.RequestID, middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(PrometheusMiddleware)
	router.Use(AuthMiddleware(domain))
}
func StartPrometheusHandlers() http.Handler {
	router := chi.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	return router
}
