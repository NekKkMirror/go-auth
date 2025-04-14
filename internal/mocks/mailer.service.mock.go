package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockMailer represents a mock mail service
type MockMailer struct {
	mock.Mock
}

func (m *MockMailer) Send(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}
