package events

import (
	"github.com/google/uuid"
	"time"
)

type EventDB struct {
	ID         uuid.UUID `db:"id"`
	Type       string    `db:"event_type"`
	ReservedTo time.Time `db:"reserved_to"`
	Payload    string    `db:"payload"`
}
