package handler

import (
	"encoding/base64"
	dto "github.com/NekKkMirror/go-auth/internal/dto/response"
	"github.com/NekKkMirror/go-auth/internal/service"

	"github.com/NekKkMirror/go-auth/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles HTTP requests related to user authentication.
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// GenerateTokens is fiber handler to issue JWT Access and Refresh tokens.
// Endpoint: POST /auth/tokens?user_id=<uuid>
func (ah *AuthHandler) GenerateTokens(c *fiber.Ctx) error {
	userID, err := utils.ParseUUIDParam(c, "user_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	clientIP := c.IP()

	access, refresh, err := ah.authService.GenerateTokenPair(c.Context(), userID, clientIP)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate tokens")
	}

	return c.Status(201).JSON(dto.NewTokenResponseDto(access, refresh))
}

// RefreshTokens is a handler to refresh Access and Refresh token pair securely.
// Endpoint: POST /auth/refresh
func (ah *AuthHandler) RefreshTokens(c *fiber.Ctx) error {
	accessToken := c.Get("Authorization")
	refreshTokenEncoded := c.Get("X-Refresh-Token")

	if accessToken == "" || refreshTokenEncoded == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing tokens in headers")
	}

	claims, err := ah.authService.VerifyJWT(accessToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid JWT token")
	}

	refreshTokenBytes, err := base64.StdEncoding.DecodeString(refreshTokenEncoded)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid refresh token encoding")
	}

	clientIP := c.IP()
	if _, err := ah.authService.ValidateRefreshToken(c.Context(), claims.ID, string(refreshTokenBytes), clientIP); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Token invalid or expired")
	}

	access, refresh, err := ah.authService.GenerateTokenPair(c.Context(), claims.UserID, clientIP)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to regenerate tokens")
	}

	return c.Status(201).JSON(dto.NewTokenResponseDto(access, refresh))
}

// RegisterAuthRoutes configures auth related routes.
func RegisterAuthRoutes(router fiber.Router, handler *AuthHandler) {
	router.Post("/auth/tokens", handler.GenerateTokens)
	router.Post("/auth/refresh", handler.RefreshTokens)
}
