package service

import "github.com/google/uuid"

// JWTServiceInterface interface for JWTService
type JWTServiceInterface interface {
	Generate(userID uuid.UUID, clientIP string) (string, string, error)
	Verify(tokenStr string) (*JWTClaims, error)
}
