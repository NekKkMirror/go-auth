package mocks

import (
	service "github.com/NekKkMirror/go-auth/internal/service/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockJWTService represents a mock of JWT service
type MockJWTService struct {
	mock.Mock
}

func (j *MockJWTService) Generate(userID uuid.UUID, clientIP string) (string, string, error) {
	args := j.Called(userID, clientIP)
	return args.String(0), args.String(1), args.Error(2)
}

func (j *MockJWTService) Verify(tokenStr string) (*service.JWTClaims, error) {
	args := j.Called(tokenStr)
	return args.Get(0).(*service.JWTClaims), args.Error(1)
}
