package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
	"github.com/xclamation/go-auth-service/internal/database"
	"github.com/xclamation/go-auth-service/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db *database.Queries
}

func NewAuthHandler(db *database.Queries) *AuthHandler {
	return &AuthHandler{db: db}
}

func generateRefreshToken() string {
	//Generate a secure random refresh token
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(token)
}

// sendEmailWarning sends an email warning to the user
func sendEmailWarning(email, newIP string) {
	logrus.Infof("Sending email warning to %s about new IP address: %s", email, newIP)
}

func (h *AuthHandler) GenerateTokenPair(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req struct {
		UserID pgtype.UUID `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.WithError(err).Error("Failed to decode requset body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Ensure the user exists in the users table
	_, err := h.db.GetUserByID(r.Context(), req.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Create the user if not found
			_, err = h.db.CreateUser(r.Context(), database.CreateUserParams{
				ID:           req.UserID,
				Email:        "user@example.com",
				PasswordHash: "test_password_hash",
			})
			if err != nil {
				logrus.WithError(err).Error("Failed to create user")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			logrus.WithError(err).Error("Failed to get user by ID")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Generate Access Token
	accessToken, err := jwt.GenerateJWT(req.UserID, r.RemoteAddr)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate access token")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate Refresh Token
	refreshToken := generateRefreshToken()
	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate refresh token hash")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store Refresh Token in the database
	_, err = h.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    req.UserID,
		TokenHash: string(hashedRefreshToken),
	})
	if err != nil {
		logrus.WithError(err).Error("Failed to store refresh token in the database")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return tokens
	json.NewEncoder(w).Encode(map[string]string{
		"acces_token":   accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) RefreshTokenPair(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.WithError(err).Error("Failed to decode request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate Refresh Token
	refreshToken, err := h.db.GetRefreshTokenByHash(r.Context(), req.RefreshToken)
	if err != nil {
		logrus.WithError(err).Error("Failed to get refresh token by hash")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Ensure the user exists in the users table
	user, err := h.db.GetUserByID(r.Context(), refreshToken.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logrus.WithError(err).Error("User not found")
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			logrus.WithError(err).Error("Failed to get user by ID")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Check if the IP address has changed
	claims, err := jwt.ValidateJWT(refreshToken.TokenHash)
	if err != nil {
		logrus.WithError(err).Error("Failed to validate access token")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if claims.IP != r.RemoteAddr {
		// Send email warning to the user
		sendEmailWarning(user.Email, r.RemoteAddr)
	}

	// Generate new Access Token
	accessToken, err := jwt.GenerateJWT(refreshToken.UserID, req.RefreshToken)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate access token")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return new Access Token
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}
