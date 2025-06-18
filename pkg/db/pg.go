package db

import (
	"context"
	"fmt"
	"github.com/Sanchir01/currency-wallet/internal/config"
	"github.com/Sanchir01/currency-wallet/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"
)

func PGXNew(cfg *config.Config, ctx context.Context) (*pgxpool.Pool, error) {
	var dsn string
	passwordpg := os.Getenv("DB_PASSWORD_PROD")
	fmt.Println(passwordpg)
	switch cfg.Env {
	case "development":
		dsn = fmt.Sprintf(
			"postgres://postgres:postgres@localhost:5441/currency-wallet?sslmode=disable",
		)
	case "production":
		dsn = fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s",
			cfg.DB.User, passwordpg,
			cfg.DB.Host, cfg.DB.Port, cfg.DB.Database,
		)
	}
	var pool *pgxpool.Pool
	var err error

	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		var err error
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}
		return nil
	}, 1, 1*time.Second)

	if err != nil {
		return nil, err
	}
	var test int
	if err := pool.QueryRow(ctx, "SELECT 1").Scan(&test); err != nil {
		log.Printf("Ошибка подключения к БД: %v", err)
		return nil, err
	}
	return pool, nil
}
