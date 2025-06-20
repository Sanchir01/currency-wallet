package wallet

type BalanceDB struct {
	Balances map[string]float32 `db:"balance"`
}

type DepositOrWithdrawRequest struct {
	Amount   float32 `json:"amount" validate:"required"`
	Currency string  `json:"currency" validate:"required"`
}
