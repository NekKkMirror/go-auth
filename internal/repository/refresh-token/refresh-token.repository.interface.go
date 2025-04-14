package repository

import (
	"context"
	"github.com/NekKkMirror/go-auth/internal/model"
)

// RefreshTokenRepositoryInterface interface for RefreshTokenRepo
type RefreshTokenRepositoryInterface interface {
	Save(ctx context.Context, refreshToken *model.RefreshToken) error
	Get(ctx context.Context, jti string) (*model.RefreshToken, error)
	Delete(ctx context.Context, jti string) error
}
