package repository

import (
	"context"
	"github.com/NekKkMirror/go-auth/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RefreshToken handles CRUD operations for refresh tokens
type RefreshToken struct {
	db *pgxpool.Pool
}

// NewRefreshTokenRepository returns an instance of AuthRepository
func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshToken {
	return &RefreshToken{db: db}
}

// Save inserts hashed token details into the database
func (ar *RefreshToken) Save(ctx context.Context, rt *model.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, associated_jti, client_ip, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`
	_, err := ar.db.Exec(ctx, query, rt.UserID, rt.TokenHash, rt.AssociatedJTI, rt.ClientIP)
	return err
}

// Get retrieves stored token details using associated JWT ID
func (ar *RefreshToken) Get(ctx context.Context, associatedJTI string) (*model.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, associated_jti, client_ip, created_at
		FROM refresh_tokens
		WHERE associated_jti = $1
		LIMIT 1
	`
	var token model.RefreshToken
	err := ar.db.QueryRow(ctx, query, associatedJTI).Scan(
		&token.ID, &token.UserID, &token.TokenHash,
		&token.AssociatedJTI, &token.ClientIP, &token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// Delete removes used refresh tokens from storage to prevent re-usage
func (ar *RefreshToken) Delete(ctx context.Context, associatedJTI string) error {
	query := `DELETE FROM refresh_tokens WHERE associated_jti = $1`
	_, err := ar.db.Exec(ctx, query, associatedJTI)
	return err
}
