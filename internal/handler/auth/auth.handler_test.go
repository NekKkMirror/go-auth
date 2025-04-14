package handler_test

import (
	"encoding/base64"
	"encoding/json"
	handler "github.com/NekKkMirror/go-auth/internal/handler/auth"
	"github.com/NekKkMirror/go-auth/internal/mocks"
	"github.com/NekKkMirror/go-auth/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http/httptest"
	"testing"
)

func setupTest() (*mocks.MockRefreshTokenRepo, *mocks.MockJWTService, *mocks.MockMailer, *fiber.App, *handler.AuthHandler) {
	mockRefreshTokenRepo := new(mocks.MockRefreshTokenRepo)
	mockJWTService := new(mocks.MockJWTService)
	mockMailer := new(mocks.MockMailer)

	authService := service.NewAuthService(mockRefreshTokenRepo, mockJWTService, mockMailer)
	authHandler := handler.NewAuthHandler(authService)

	app := fiber.New()

	return mockRefreshTokenRepo, mockJWTService, mockMailer, app, authHandler
}

func TestAuthHandler_Tokens(t *testing.T) {
	mockRefreshTokenRepo, mockJWTService, _, app, authHandler := setupTest()

	t.Run("Successful Token Generation", func(t *testing.T) {
		userID := uuid.New()
		clientIP := "0.0.0.0"
		expectedAccessToken := "mockAccessToken"

		mockJWTService.On("Generate", userID, clientIP).Return(expectedAccessToken, "jti-mock", nil)
		mockRefreshTokenRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

		app.Post("/auth/tokens", func(c *fiber.Ctx) error {
			return authHandler.GenerateTokens(c)
		})

		req := httptest.NewRequest("POST", "/auth/tokens?user_id="+userID.String(), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Forwarded-For", clientIP)

		resp, err := app.Test(req, -1)
		assert.Nil(t, err)

		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, 201, resp.StatusCode)

		var responseData map[string]string
		err = json.Unmarshal(body, &responseData)
		assert.Nil(t, err)
		assert.Equal(t, expectedAccessToken, responseData["access_token"])

		decodedToken, err := base64.StdEncoding.DecodeString(responseData["refresh_token"])
		assert.Nil(t, err)
		assert.NotEmpty(t, decodedToken)

		mockJWTService.AssertCalled(t, "Generate", userID, clientIP)
		mockRefreshTokenRepo.AssertCalled(t, "Save", mock.Anything, mock.Anything)
	})
}
