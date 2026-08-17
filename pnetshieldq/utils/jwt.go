package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"time"
)

func jwtSigningKey() ([]byte, error) {
	key := os.Getenv("JWT_SIGNING_KEY")
	if key == "" {
		return nil, errors.New("JWT_SIGNING_KEY is not configured")
	}
	return []byte(key), nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	key, err := jwtSigningKey()
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	return token, err
}

func GetUserIDFromJWT(r *http.Request) int {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return -1
	}
	tokenString := authHeader[7:]
	token, err := ParseToken(tokenString)
	if err != nil {
		return -1
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return int(claims["user_id"].(float64))
	}
	return -1
}

func GenerateJWT(userID int, username string, password string) (string, error) {
	key, err := jwtSigningKey()
	if err != nil {
		return "", err
	}
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(key)
	return tokenString, err
}
func GetUserIDFromTokenString(tokenString string) (int, error) {
	token, err := ParseToken(tokenString)
	if err != nil {
		return -1, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return int(claims["user_id"].(float64)), nil
	}
	return -1, nil
}
