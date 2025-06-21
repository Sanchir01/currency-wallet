package wallet

import (
	"github.com/Sanchir01/currency-wallet/internal/domain/models"
	"github.com/Sanchir01/currency-wallet/pkg/api"
	"github.com/google/uuid"
)

type BalanceDB struct {
	Balances map[string]float32 `db:"balance"`
}
type CurrencyWalletResponse struct {
	Rates map[string]float32
}
type DepositOrWithdrawRequest struct {
	Amount   float32 `json:"amount" validate:"required"`
	Currency string  `json:"currency" validate:"required"`
}

type DepositOrWithdrawResponse struct {
	api.Response
	models.CurrencyWallet
}

type ExchangeRequest struct {
	FromCurrency string  `json:"from_currency" validate:"required"`
	ToCurrency   string  `json:"to_currency" validate:"required"`
	Amount       float32 `json:"amount" validate:"required"`
}
type ExchangeResponse struct {
	api.Response
	Message         string             `json:"message"`
	ExchangedAmount float32            `json:"exchanged_amount"`
	NewBalance      map[string]float32 `json:"new_balance"`
}

type KafkaPayloadNotification struct {
	Amount       float32            `json:"amount"`
	UserId       uuid.UUID          `json:"user_id"`
	BalanceAfter map[string]float32 `json:"balance_after"`
}
type ExchangeRateToCurrency struct {
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Rate         float32 `json:"rate"`
}
type DefaultRedis struct {
	Rates map[string]float64 `json:"rates"`
}
