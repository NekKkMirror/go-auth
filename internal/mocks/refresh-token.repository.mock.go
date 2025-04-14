package mocks

import (
	"context"

	"github.com/NekKkMirror/go-auth/internal/model"
	"github.com/stretchr/testify/mock"
)

// MockRefreshTokenRepo presents a mock of the RefreshToken repository
type MockRefreshTokenRepo struct {
	mock.Mock
}

func (m *MockRefreshTokenRepo) Save(ctx context.Context, refreshToken *model.RefreshToken) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func (m *MockRefreshTokenRepo) Get(ctx context.Context, jti string) (*model.RefreshToken, error) {
	args := m.Called(ctx, jti)
	return args.Get(0).(*model.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepo) Delete(ctx context.Context, jti string) error {
	args := m.Called(ctx, jti)
	return args.Error(0)
}
