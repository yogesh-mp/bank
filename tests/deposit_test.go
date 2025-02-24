package tests

import (
	"banking-ledger-service/internal/handlers"
	"banking-ledger-service/internal/models"
	"banking-ledger-service/tests/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeposit_Success(t *testing.T) {
	mockDB := new(mocks.MockDB)
	mockQueue := new(mocks.MockQueue)

	mockDB.On("GetAccount", mock.Anything).Return(&models.Account{ID: 1, Name: "John Doe", Balance: 1000.00}, nil)
	mockQueue.On("PublishMessage", mock.Anything).Return(nil)

	transaction := models.Transaction{AccountID: 1, Amount: 500.00}
	reqBody, _ := json.Marshal(transaction)

	req := httptest.NewRequest("POST", "/deposit", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handlers.Deposit(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Deposit request sent to queue")

	mockDB.AssertExpectations(t)
	mockQueue.AssertExpectations(t)
}

func TestDeposit_NegativeAmount(t *testing.T) {
	transaction := models.Transaction{AccountID: 1, Amount: -500.00}
	reqBody, _ := json.Marshal(transaction)

	req := httptest.NewRequest("POST", "/deposit", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handlers.Deposit(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Deposit amount must be greater than zero")
}

func TestDeposit_ZeroAmount(t *testing.T) {
	transaction := models.Transaction{AccountID: 1, Amount: 0.00}
	reqBody, _ := json.Marshal(transaction)

	req := httptest.NewRequest("POST", "/deposit", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handlers.Deposit(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Deposit amount must be greater than zero")
}

func TestDeposit_NonExistentAccount(t *testing.T) {
	mockDB := new(mocks.MockDB)

	mockDB.On("GetAccount", mock.Anything).Return(nil, errors.New("account not found"))

	transaction := models.Transaction{AccountID: 99, Amount: 500.00}
	reqBody, _ := json.Marshal(transaction)

	req := httptest.NewRequest("POST", "/deposit", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handlers.Deposit(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Account not found")

	mockDB.AssertExpectations(t)
}

func TestDeposit_QueueFailure(t *testing.T) {
	mockDB := new(mocks.MockDB)
	mockQueue := new(mocks.MockQueue)

	mockDB.On("GetAccount", mock.Anything).Return(&models.Account{ID: 1, Name: "John Doe", Balance: 1000.00}, nil)
	mockQueue.On("PublishMessage", mock.Anything).Return(errors.New("queue failure"))

	transaction := models.Transaction{AccountID: 1, Amount: 500.00}
	reqBody, _ := json.Marshal(transaction)

	req := httptest.NewRequest("POST", "/deposit", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handlers.Deposit(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to queue deposit transaction")

	mockDB.AssertExpectations(t)
	mockQueue.AssertExpectations(t)
}
