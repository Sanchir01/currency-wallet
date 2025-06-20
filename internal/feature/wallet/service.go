package wallet

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Sanchir01/currency-wallet/internal/domain/contants"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"

	"github.com/Sanchir01/currency-wallet/internal/domain/models"
	walletsv1 "github.com/Sanchir01/wallets-proto/gen/go/wallets"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ServiceWallets interface {
	Balance(ctx context.Context, id uuid.UUID) (*models.CurrencyWallet, error)
	DepositOrWithdrawBalance(ctx context.Context, id uuid.UUID, amount float32, currency string, tx pgx.Tx, typedepo contextkey.OperationType) (*models.CurrencyWallet, error)
}
type Service struct {
	repository ServiceWallets
	exchanger  walletsv1.ExchangeServiceClient
	redisdb    *redis.Client
	primaryDB  *pgxpool.Pool
	log        *slog.Logger
}

func NewService(r ServiceWallets, primaryDB *pgxpool.Pool, redisdb *redis.Client, exchanger walletsv1.ExchangeServiceClient, log *slog.Logger) *Service {
	return &Service{
		repository: r,
		exchanger:  exchanger,
		redisdb:    redisdb,
		log:        log,
		primaryDB:  primaryDB,
	}
}

type DefaultRedis struct {
	Rates map[string]float64 `json:"rates"`
}

func (s *Service) GetCurrencyWallets(ctx context.Context) (*models.CurrencyWallet, error) {
	const op = "Wallet.Service.GetCurrencyWallets"
	log := s.log.With(slog.String("op", op))

	dataredis, err := s.redisdb.Get(ctx, contextkey.ExchangerCurrencyCtxKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			log.Info("cache miss (key not found), fetching from exchanger")
		} else {
			log.Error("redis error", slog.String("error", err.Error()))
		}

		exchangerdata, err := s.exchanger.GetExchangeRates(ctx, &emptypb.Empty{})
		if err != nil {
			log.Error("failed to get exchange rates", slog.String("error", err.Error()))
			return nil, err
		}

		log.Debug("received exchange rates", slog.Any("rates", exchangerdata.Rates))

		data, err := json.Marshal(exchangerdata.Rates)
		if err != nil {
			log.Error("failed to marshal rates", slog.String("error", err.Error()))
			return nil, err
		}

		if err := s.redisdb.Set(ctx, contextkey.ExchangerCurrencyCtxKey, data, 5*time.Second).Err(); err != nil {
			log.Error("failed to set data to redis", slog.String("error", err.Error()))
			return nil, err
		}

		return &models.CurrencyWallet{Balances: exchangerdata.Rates}, nil
	}

	var rates DefaultRedis
	if err := json.Unmarshal(dataredis, &rates); err != nil {
		log.Error("failed to unmarshal redis data (float64)", slog.String("error", err.Error()))
		return nil, err
	}

	balances := make(map[string]float32, len(rates.Rates))
	for k, v := range rates.Rates {
		balances[k] = float32(v)
	}

	log.Info("successfully unmarshaled redis data", slog.Any("rates", balances))
	return &models.CurrencyWallet{Balances: balances}, nil
}

func (s *Service) GetBalance(ctx context.Context, id uuid.UUID) (*models.CurrencyWallet, error) {
	const op = "Wallet.Service.GetBalance"
	log := s.log.With(slog.String("op", op))
	data, err := s.repository.Balance(ctx, id)
	if err != nil {
		log.Error("failed to get balance", slog.String("error", err.Error()))
		return nil, err
	}
	log.Info("successfully get balance")
	return data, nil
}

func (s *Service) WalletDepositOrWithDraw(ctx context.Context, id uuid.UUID, currency string, amount float32, typedepo contextkey.OperationType) (*models.CurrencyWallet, error) {
	const op = "Wallet.Service.GetBalance"
	log := s.log.With(slog.String("op", op))
	log.Warn("DepositOrWithdrawBalance props", amount)
	conn, err := s.primaryDB.Acquire(ctx)
	if err != nil {

		return nil, err
	}
	defer conn.Release()
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Error("tx error", err.Error())
		return nil, err
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
				log.Error("rollback error", rollbackErr.Error())
				return
			}
		}

	}()
	data, err := s.repository.DepositOrWithdrawBalance(ctx, id, amount, currency, tx, typedepo)
	if err != nil {
		log.Error("failed to deposit balance", slog.String("error", err.Error()))
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		log.Error("failed to commit transaction", slog.String("error", err.Error()))
		return nil, err
	}
	log.Info("getting balance for currency")
	return data, nil
}
