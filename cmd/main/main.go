package main

import (
	"context"
	"errors"
	"github.com/Sanchir01/currency-wallet/internal/app"
	httphandlers "github.com/Sanchir01/currency-wallet/internal/http"
	httpserver "github.com/Sanchir01/currency-wallet/pkg/server/http"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// @title 🚀 Currency Wallet
// @version         1.0
// @description This is a sample server celler
// @termsOfService  http://swagger.io/terms/

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @contact.name GitHub
// @contact.url https://github.com/Sanchir01
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	env, err := app.NewApp(ctx)
	if err != nil {
		log.Fatal(err)
	}
	env.Lg.Info("Currency wallet initialized")
	serve := httpserver.NewHTTPServer(env.Cfg.HTTPServer.Host, env.Cfg.HTTPServer.Port,
		env.Cfg.HTTPServer.Timeout, env.Cfg.HTTPServer.IdleTimeout)
	prometheusserver := httpserver.NewHTTPServer(env.Cfg.Prometheus.Host, env.Cfg.Prometheus.Port, env.Cfg.Prometheus.Timeout,
		env.Cfg.Prometheus.IdleTimeout)

	go func() {
		if err := serve.Run(httphandlers.StartHTTTPHandlers(env.Handlers, env.Cfg.Domain)); err != nil {
			if !errors.Is(err, context.Canceled) {
				env.Lg.Error("Listen server error", slog.String("error", err.Error()))
				return
			}
		}
	}()
	go func() {
		if err := prometheusserver.Run(httphandlers.StartPrometheusHandlers()); err != nil {
			if !errors.Is(err, context.Canceled) {
				env.Lg.Error("Listen prometheus server error", slog.String("error", err.Error()))
				return
			}
		}
	}()
	<-ctx.Done()

	if err := serve.Gracefull(ctx); err != nil {
		env.Lg.Error("server gracefull")
	}
	if err := env.DB.Close(); err != nil {
		env.Lg.Error("Close database", slog.String("error", err.Error()))
	}
}
