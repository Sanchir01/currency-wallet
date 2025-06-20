package contextkey

type ContextKey string

const UserIDCtxKey ContextKey = "userID"

const ExchangerCurrencyCtxKey string = "exchangerCurrency"

var Currencies = []string{"USD", "EUR", "RUB"}

type OperationType string

const (
	OperationTypeDeposit  OperationType = "DEPOSIT"
	OperationTypeWithdraw OperationType = "WITHDRAW"
	OperationTypeTransfer OperationType = "TRANSFER"
)
