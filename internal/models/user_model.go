package models

import (
	"time"

	"github.com/google/uuid"
)

type UserID = uuid.UUID

type User struct {
	UserID       UserID    `json:"user_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}
