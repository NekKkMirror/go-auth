package model

import (
	"github.com/google/uuid"
	"time"
)

// RefreshToken represents stored hashed refresh token information
type RefreshToken struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	TokenHash     string
	AssociatedJTI string
	ClientIP      string
	CreatedAt     time.Time
}
