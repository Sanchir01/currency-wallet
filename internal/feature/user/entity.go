package user

import (
	"github.com/google/uuid"
	"time"
)

type DatabaseUser struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	Name      string    `db:"username"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Version   int64     `db:"version"`
}

type GetAllUserCoinsInfo struct {
	Coins        int64        `json:"coins"`
	Inventory    []Inventory  `json:"inventory"`
	CoinsHistory CoinsHistory `json:"coinsHistory"`
}
type Inventory struct {
	Type     string `json:"type"`
	Quantity int64  `json:"quantity"`
}
type CoinsHistory struct {
	Received []SendInfo `json:"received"`
	Send     []SendInfo `json:"send"`
}
type SendInfo struct {
	Username string `json:"fromUser"`
	Amount   int64  `json:"amount"`
}
type ProductCount struct {
	Title string `db:"title"`
	Count int64  `db:"quantity"`
}
