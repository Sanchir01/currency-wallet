package wallet

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	contextkey "github.com/Sanchir01/currency-wallet/internal/domain/contants"
	"github.com/Sanchir01/currency-wallet/internal/domain/models"
	"github.com/Sanchir01/currency-wallet/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type Repository struct {
	primaryDB *pgxpool.Pool
	log       *slog.Logger
}

func NewRepository(primaryDB *pgxpool.Pool, log *slog.Logger) *Repository {
	return &Repository{
		primaryDB: primaryDB,
		log:       log,
	}
}

func (r *Repository) Balance(ctx context.Context, id uuid.UUID) (*models.CurrencyWallet, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	query, args, err := sq.Select("balance, currency").
		From("wallets").
		Where(sq.Eq{"user_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, utils.ErrorQueryString
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	balances := make(map[string]float32)

	for rows.Next() {
		var balance float32
		var currency string
		if err := rows.Scan(&balance, &currency); err != nil {
			return nil, err
		}

		balances[currency] = balance
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return (*models.CurrencyWallet)(&BalanceDB{
		Balances: balances,
	}), nil
}

func (r *Repository) CreateManyWallets(ctx context.Context, userID uuid.UUID, tx pgx.Tx) error {
	builder := sq.Insert("wallets").
		Columns("user_id", "currency", "balance").
		PlaceholderFormat(sq.Dollar)

	for _, currency := range contextkey.Currencies {
		builder = builder.Values(userID, currency, 0.0)
	}
	query, args, err := builder.ToSql()
	if err != nil {
		return utils.ErrorQueryString
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		slog.Error("failed to insert wallets", slog.Any("err", err))
		return err
	}
	return nil
}

func (r *Repository) DepositOrWithdrawBalance(
	ctx context.Context,
	id uuid.UUID,
	amount float32,
	currency string,
	tx pgx.Tx,
	typedepo contextkey.OperationType,
) (*models.CurrencyWalletDB, error) {
	var query string
	r.log.Warn("depo props", amount)

	if typedepo == contextkey.OperationTypeWithdraw {
		query = `
		WITH updated_wallet AS (
			UPDATE wallets
			SET balance = balance - $1
			WHERE user_id = $2 AND currency = $3 AND balance >= $1
			RETURNING id, user_id
		)
		SELECT w.id, w.currency, w.balance
		FROM wallets w
		JOIN updated_wallet uw ON w.user_id = uw.user_id;
		`
	} else {
		query = `
		WITH updated_wallet AS (
			UPDATE wallets
			SET balance = balance + $1
			WHERE user_id = $2 AND currency = $3
			RETURNING id, user_id
		)
		SELECT w.id, w.currency, w.balance
		FROM wallets w
		JOIN updated_wallet uw ON w.user_id = uw.user_id;
		`
	}

	args := []interface{}{amount, id, currency}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	balances := make(map[string]float32)
	var walletID uuid.UUID
	first := true

	for rows.Next() {
		var currency string
		var balance float32
		var currentWalletID uuid.UUID

		if err := rows.Scan(&currentWalletID, &currency, &balance); err != nil {
			return nil, err
		}
		balances[currency] = balance

		if first {
			walletID = currentWalletID
			first = false
		}
	}

	if len(balances) == 0 {
		return nil, errors.New("insufficient funds or wallet not found")
	}

	return &models.CurrencyWalletDB{
		WalletID: walletID,
		CurrencyWallet: models.CurrencyWallet{
			Balances: balances,
		},
	}, nil
}

func (r *Repository) SetTransaction(ctx context.Context, walletID uuid.UUID,
	amount float32,
	typetransaction contextkey.OperationType,
	senderID *uuid.UUID,
	tx pgx.Tx) error {
	queryBuilder := sq.Insert("transactions").
		Columns("wallet_id", "amount", "type").
		Values(walletID, amount, typetransaction).
		PlaceholderFormat(sq.Dollar)

	if senderID != nil {
		queryBuilder = sq.Insert("transactions").
			Columns("wallet_id", "amount", "type", "sender_wallet_id").
			Values(walletID, amount, typetransaction, *senderID).
			PlaceholderFormat(sq.Dollar)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return utils.ErrorQueryString
	}

	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
