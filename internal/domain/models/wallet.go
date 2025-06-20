package models

import "github.com/google/uuid"

type CurrencyWallet struct {
	Balances map[string]float32
}

type CurrencyWalletDB struct {
	CurrencyWallet
	WalletID uuid.UUID
}
