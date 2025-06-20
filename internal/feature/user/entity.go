package user

import (
	"github.com/Sanchir01/currency-wallet/pkg/api"
	"github.com/google/uuid"
	"time"
)

type DatabaseUser struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	Name      string    `db:"username"`
	Password  []byte    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Version   int64     `db:"version"`
}
type AuthRequest struct {
	Email    string `json:"email" validate:"required"`
	Username string `json:"username" validate:"required,min=1,max=100"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthResponse struct {
	api.Response
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	api.Response
	Email    string `json:"email"`
	Username string `json:"username" `
}
