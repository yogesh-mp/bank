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

func TestCreateAccount_Success(t *testing.T) {
	mockDB := new(mocks.MockDB)
	mockQueue := new(mocks.MockQueue)

	// Mock storage method
	mockDB.On("GetAccount", mock.Anything).Return(nil, nil)                 // Account does not exist
	mockDB.On("CreateAccount", mock.Anything, mock.Anything).Return(1, nil) // Account creation succeeds

	// Mock queue
	mockQueue.On("PublishMessage", mock.Anything).Return(nil)

	// Prepare request
	account := models.Account{Name: "John Doe", Balance: 1000.00}
	reqBody, _ := json.Marshal(account)

	req := httptest.NewRequest("POST", "/create-account", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handlers.CreateAccount(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	mockDB.AssertExpectations(t)
	mockQueue.AssertExpectations(t)
}

func TestCreateAccount_AlreadyExists(t *testing.T) {
	mockDB := new(mocks.MockDB)

	// Mock account already exists
	mockDB.On("GetAccount", mock.Anything).Return(&models.Account{ID: 1, Name: "John Doe", Balance: 500}, nil)

	reqBody, _ := json.Marshal(models.Account{Name: "John Doe", Balance: 1000.00})
	req := httptest.NewRequest("POST", "/create-account", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handlers.CreateAccount(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
	assert.Contains(t, rec.Body.String(), "Account already exists")

	mockDB.AssertExpectations(t)
}

func TestCreateAccount_InvalidPayload(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/create-account", bytes.NewReader([]byte("{invalid json}")))
	rec := httptest.NewRecorder()

	handlers.CreateAccount(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid request")
}

func TestCreateAccount_QueueFailure(t *testing.T) {
	mockDB := new(mocks.MockDB)
	mockQueue := new(mocks.MockQueue)

	mockDB.On("GetAccount", mock.Anything).Return(nil, nil)                 // Account does not exist
	mockDB.On("CreateAccount", mock.Anything, mock.Anything).Return(1, nil) // Account creation succeeds

	// Simulate queue failure
	mockQueue.On("PublishMessage", mock.Anything).Return(errors.New("queue failure"))

	reqBody, _ := json.Marshal(models.Account{Name: "John Doe", Balance: 1000.00})
	req := httptest.NewRequest("POST", "/create-account", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handlers.CreateAccount(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to queue account creation")

	mockDB.AssertExpectations(t)
	mockQueue.AssertExpectations(t)
}
