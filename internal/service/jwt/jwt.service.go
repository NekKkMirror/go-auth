package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTService structure with JWT signature key
type JWTService struct {
	signingKey []byte
}

// JWTClaims define the structure of data embedded in the JWT
type JWTClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	ClientIP string    `json:"client_ip"`
	jwt.RegisteredClaims
}

// NewJWTService service constructor that accepts a key from the environment
func NewJWTService(signingKey string) *JWTService {
	return &JWTService{
		signingKey: []byte(signingKey),
	}
}

// Generate generates secure JWT tokens including client info and JTI identifier
func (j *JWTService) Generate(userID uuid.UUID, clientIP string) (string, string, error) {
	jti := uuid.New().String()

	claims := JWTClaims{
		UserID:   userID,
		ClientIP: clientIP,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Issuer:    "go-auth-service",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &claims)
	signedToken, err := token.SignedString(j.signingKey)

	return signedToken, jti, err
}

// Verify validates JWT token and extracts payload securely
func (j *JWTService) Verify(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.signingKey, nil
	})

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
