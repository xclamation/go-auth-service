package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
)

var jwtKey []byte

func init() {
	// Load environment variable from .env file
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// Get the JWT secret from the environment variable
	jwtKey = []byte(os.Getenv("JWT_SECRET"))
	if jwtKey == nil {
		panic("JWT_SECRET environment variable is not set")
	}
}

type Claims struct {
	UserID pgtype.UUID `json:"user_id"`
	IP     string      `json:"ip"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID pgtype.UUID, ip string) (string, error) {
	exprirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		UserID: userID,
		IP:     ip,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exprirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return claims, nil
}
