package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/ctu-ikz/timetable-be/models"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .en file")
	}
}

func parseDuration(value string) time.Duration {
	d, err := time.ParseDuration(value + "s")
	if err != nil {
		return 0
	}
	return d
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(password string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func CreateTokens(user *models.UserDB) (string, string, error) {
	if user.ID == nil {
		return "", "", fmt.Errorf("user ID cannot be nil")
	}

	accessTokenExpiry := parseDuration(os.Getenv("ACCESS_TOKEN_EXPIRY"))
	refreshTokenExpiry := parseDuration(os.Getenv("REFRESH_TOKEN_EXPIRY"))
	accessSecret := os.Getenv("ACCESS_SECRET")
	refreshSecret := os.Getenv("REFRESH_SECRET")

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       *user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(accessTokenExpiry).Unix(),
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  *user.ID,
		"exp": time.Now().Add(refreshTokenExpiry).Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(accessSecret))
	if err != nil {
		return "", "", err
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(refreshSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func VerifyToken(tokenString string, isRefresh bool) error {
	var secret []byte

	if isRefresh {
		secret = []byte(os.Getenv("REFRESH_SECRET"))
	} else {
		secret = []byte(os.Getenv("ACCESS_SECRET"))
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return fmt.Errorf("invalid token: missing expiration")
	}
	if time.Now().Unix() > int64(exp) {
		return fmt.Errorf("token has expired")
	}

	return nil
}

func GetUserIDFromToken(tokenString string, isRefresh bool) (string, error) {
	var secret []byte
	if isRefresh {
		secret = []byte(os.Getenv("REFRESH_SECRET"))
	} else {
		secret = []byte(os.Getenv("ACCESS_SECRET"))
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if id, ok := claims["id"].(string); ok {
			return id, nil
		}
		if id, ok := claims["id"].(float64); ok {
			return fmt.Sprintf("%.0f", id), nil
		}
	}
	return "", fmt.Errorf("id not found in token claims")
}
