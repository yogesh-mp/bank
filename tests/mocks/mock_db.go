package mocks

import (
	"banking-ledger-service/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockDB simulates the behavior of storage.DB methods
type MockDB struct {
	mock.Mock
}

// Mock CreateAccount method
func (m *MockDB) CreateAccount(name string, balance float64) (int, error) {
	args := m.Called(name, balance)
	return args.Int(0), args.Error(1)
}

// Mock GetAccount method
func (m *MockDB) GetAccount(id int) (*models.Account, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

// Mock UpdateBalance method
func (m *MockDB) UpdateBalance(accountID int, amount float64, operation string) error {
	args := m.Called(accountID, amount, operation)
	return args.Error(0)
}

// Mock AddTransaction method
func (m *MockDB) AddTransaction(accountID int, amount float64, txType string) error {
	args := m.Called(accountID, amount, txType)
	return args.Error(0)
}
