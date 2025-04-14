package service

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/NekKkMirror/go-auth/internal/model"
	repository "github.com/NekKkMirror/go-auth/internal/repository/refresh-token"
	service "github.com/NekKkMirror/go-auth/internal/service/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
)

// AuthService provides methods related to user authentication logic.
type AuthService struct {
	refreshTokenRepo repository.RefreshTokenRepositoryInterface
	jwtService       service.JWTServiceInterface
	mailerService    Mailer
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(refreshTokenRepo repository.RefreshTokenRepositoryInterface, jwtService service.JWTServiceInterface, mailerService Mailer) *AuthService {
	return &AuthService{
		refreshTokenRepo: refreshTokenRepo,
		jwtService:       jwtService,
		mailerService:    mailerService,
	}
}

// GenerateTokenPair generates JWT Access and Refresh tokens and store the refresh token securely.
func (as *AuthService) GenerateTokenPair(ctx context.Context, userID uuid.UUID, clientIP string) (access string, refresh string, err error) {
	access, jti, err := as.jwtService.Generate(userID, clientIP)
	if err != nil {
		return "", "", err
	}

	refreshRaw := uuid.New().String()
	refreshHashBytes, err := bcrypt.GenerateFromPassword([]byte(refreshRaw), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	if err := as.refreshTokenRepo.Save(ctx, &model.RefreshToken{
		UserID:        userID,
		TokenHash:     string(refreshHashBytes),
		AssociatedJTI: jti,
		ClientIP:      clientIP,
	}); err != nil {
		return "", "", err
	}

	refresh = base64.StdEncoding.EncodeToString([]byte(refreshRaw))
	return access, refresh, nil
}

// VerifyJWT validates the authenticity of a JWT token.
func (as *AuthService) VerifyJWT(tokenStr string) (*service.JWTClaims, error) {
	return as.jwtService.Verify(tokenStr)
}

// ValidateRefreshToken validates the refresh token and checks client IP changes.
func (as *AuthService) ValidateRefreshToken(ctx context.Context, jti string, refreshRaw string, clientIP string) (*model.RefreshToken, error) {
	rt, err := as.refreshTokenRepo.Get(ctx, jti)
	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(rt.TokenHash), []byte(refreshRaw)) != nil {
		return nil, ErrInvalidToken
	}

	if rt.ClientIP != clientIP {
		go as.sendEmail(clientIP, "yourname@examle.com") // мок для упрощения
	}

	if err := as.refreshTokenRepo.Delete(ctx, jti); err != nil {
		return nil, err
	}

	return rt, nil
}

// sendEmail simulates notifying user via email about IP changes (mock).
func (as *AuthService) sendEmail(newIP, recipient string) {
	subject := "IP Address Change Notification"
	body := "Your IP address has changed to: " + newIP

	if err := as.mailerService.Send(recipient, subject, body); err != nil {
		println("Failed to send email:", err.Error())
	}
}
