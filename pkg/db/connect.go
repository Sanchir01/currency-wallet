package db

import (
	"context"
	"github.com/Sanchir01/exchanger/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"os"
)

type Database struct {
	PrimaryDB *pgxpool.Pool
	RedisDB   *redis.Client
}

func NewDataBases(cfg *config.Config, ctx context.Context) (*Database, error) {
	pgxdb, err := PGXNew(cfg, ctx)
	if err != nil {
		return nil, err
	}
	redisdb, err := RedisConnect(context.Background(), cfg.RedisDB.Host, cfg.RedisDB.Port,
		os.Getenv("REDIS_PASSWORD"), cfg.Env,
		cfg.RedisDB.DBNumber, cfg.RedisDB.Retries)
	if err != nil {
		return nil, err
	}
	return &Database{PrimaryDB: pgxdb, RedisDB: redisdb}, nil
}

func (databases *Database) Close() error {
	databases.PrimaryDB.Close()
	return nil
}
