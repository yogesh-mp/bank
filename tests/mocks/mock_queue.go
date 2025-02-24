package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockQueue simulates the message queue behavior
type MockQueue struct {
	mock.Mock
}

// Mock PublishMessage method
func (m *MockQueue) PublishMessage(message string) error {
	args := m.Called(message)
	return args.Error(0)
}
