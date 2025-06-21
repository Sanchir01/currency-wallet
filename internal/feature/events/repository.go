package events

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/Sanchir01/currency-wallet/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

type Repository struct {
	primaryDB *pgxpool.Pool
}

func NewRepository(primaryDB *pgxpool.Pool) *Repository {
	return &Repository{
		primaryDB,
	}
}

func (r *Repository) CreateEvent(ctx context.Context, eventType, payload string, tx pgx.Tx) (uuid.UUID, error) {
	thistime := time.Now()
	reservedtime := thistime.Add(1 * time.Hour)
	query, args, err := sq.Insert("events").
		Columns("event_type", "payload", "reserved_to").
		Values(eventType, payload, reservedtime).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return uuid.Nil, err
	}
	var id uuid.UUID
	if err := tx.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *Repository) GetManyEvents(ctx context.Context, limit uint64) ([]*EventDB, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	query, args, err := sq.
		Select("id,event_type,payload,reserved_to").
		From("events").
		Limit(limit).
		Where(sq.Eq{"status": "new"}).
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
	events := make([]*EventDB, 0)
	for rows.Next() {
		var event EventDB
		if err := rows.Scan(&event.ID, &event.Type, &event.Payload, &event.ReservedTo); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, utils.ErrorNotFoundRows
			}
			return nil, err
		}
		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (r *Repository) SetDone(ctx context.Context, ids []uuid.UUID) error {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query, args, err := sq.Update("events").
		Set("status", "done").
		Where(sq.Eq{"id": ids}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return utils.ErrorQueryString
	}
	slog.Error("ids events", ids)
	if _, err := conn.Exec(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
