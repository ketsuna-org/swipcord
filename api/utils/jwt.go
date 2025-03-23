package utils

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:   "api",
		Subject:  "auth",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Audience: jwt.ClaimStrings{"user"},
		ID:       userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(GetEnv("JWT_SECRET", "secret")))
}

func VerifyJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetEnv("JWT_SECRET", "secret")), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", err
	}
	log.Println(claims)
	return claims["jti"].(string), nil
}
