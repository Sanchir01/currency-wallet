package wallet

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	contextkey "github.com/Sanchir01/currency-wallet/internal/contants"
	"github.com/Sanchir01/currency-wallet/internal/domain/models"
	walletsv1 "github.com/Sanchir01/wallets-proto/gen/go/wallets"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ServiceWallets interface {
}
type Service struct {
	repository ServiceWallets
	exchanger  walletsv1.ExchangeServiceClient
	redisdb    *redis.Client
	log        *slog.Logger
}

func NewService(r ServiceWallets, redisdb *redis.Client, exchanger walletsv1.ExchangeServiceClient, log *slog.Logger) *Service {
	return &Service{
		repository: r,
		exchanger:  exchanger,
		redisdb:    redisdb,
		log:        log,
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
