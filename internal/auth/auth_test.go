package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xclamation/go-auth-service/internal/database"
)

// MockDB is a mock implementation of the database.Queries interface
type MockDB struct {
	mock.Mock
	database.DBTX
}

// CreateRefreshToken mocks the CreateRefreshToken method
func (m *MockDB) CreateRefreshToken(ctx context.Context, arg database.CreateRefreshTokenParams) (database.RefreshToken, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(database.RefreshToken), args.Error(1)
}

// GetRefreshTokenByHash mocks the GetRefreshTokenByHash method
func (m *MockDB) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (database.RefreshToken, error) {
	args := m.Called(ctx, tokenHash)
	return args.Get(0).(database.RefreshToken), args.Error(1)
}

// TestGenerateTokenPair tests the GenerateTokenPair handler
func TestGenerateTokenPair(t *testing.T) {
	// Create a mock database
	mockDB := new(MockDB)
	defer mockDB.AssertExpectations(t)

	// Create a new auth handler
	authHandler := NewAuthHandler(database.New(mockDB))

	// Create a test request
	userID := uuid.New()
	pgUUID := pgtype.UUID{}
	copy(pgUUID.Bytes[:], userID[:])
	reqBody := map[string]pgtype.UUID{
		"user_id": pgUUID,
	}
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body^ %v", err)
	}

	req, err := http.NewRequest("POST", "/token", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		t.Fatalf("Failed to create request^ %v", err)
	}
	defer req.Body.Close()

	// Crate a test response recorder
	rr := httptest.NewRecorder()

	// Mock the database call
	mockDB.On("CreateRefreshToken", mock.Anything, mock.Anything).Return(database.RefreshToken{}, nil)

	// Call the handler
	authHandler.GenerateTokenPair(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	assert.NotEmpty(t, response["access_token"])
	assert.NotEmpty(t, response["refresh_tolen"])
}

// TestRefreshTokenPair tests the RefreshTokenPair handler
func TestRefreshTokenPair(t *testing.T) {
	// Create a mock database
	mockDB := new(MockDB)
	defer mockDB.AssertExpectations(t)

	// Create a new auth handler
	authHandler := NewAuthHandler(database.New(mockDB))

	// Create a test request
	refreshToken := "valid_refresh_token"
	reqBody := map[string]string{
		"refresh_token": refreshToken,
	}
	reqBodyJSON, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "/refresh", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	defer req.Body.Close()

	// Create a test response recorder
	rr := httptest.NewRecorder()

	// Mock the database call
	userID := uuid.New()
	pgUUID := pgtype.UUID{}
	copy(pgUUID.Bytes[:], userID[:])
	mockDB.On("GetRefreshTokenByHash", mock.Anything, refreshToken).Return(database.RefreshToken{
		UserID: pgUUID,
	}, nil)

	// Check the handler
	authHandler.RefreshTokenPair(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	assert.NotEmpty(t, response["access_token"])
}
