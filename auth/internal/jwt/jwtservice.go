package jwt

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	refreshSecret []byte
	accessSecret  []byte
}

func NewJWTService(refreshSecretStr string, accessSecretStr string) *JWTService {
	if refreshSecretStr == "" {
		log.Panic(fmt.Errorf("refreshSecretStr is empty"))
	}
	if accessSecretStr == "" {
		log.Panic(fmt.Errorf("accessSecretStr is empty"))
	}

	// ref: https://golang-jwt.github.io/jwt/usage/signing_methods/#signing-methods-and-key-types
	refreshSecret, err := base64.StdEncoding.DecodeString(refreshSecretStr)
	if err != nil {
		log.Panic(fmt.Errorf("couldn't convert refreshSecret into bytes: %w", err))
	}
	accessSecret, err := base64.StdEncoding.DecodeString(refreshSecretStr)
	if err != nil {
		log.Panic(fmt.Errorf("couldn't convert accessSecret into bytes: %w", err))
	}

	return &JWTService{
		refreshSecret: refreshSecret,
		accessSecret:  accessSecret,
	}
}

func (j *JWTService) NewAccessToken(data map[string]any) (string, error) {
	return j.newToken(data, j.accessSecret)
}

func (j *JWTService) NewRefreshToken(data map[string]any) (string, error) {
	return j.newToken(data, j.refreshSecret)
}

func (j *JWTService) newToken(data map[string]any, secret []byte) (string, error) {
	mapClaims := jwt.MapClaims{
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	}
	for k, v := range data {
		mapClaims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, mapClaims)
	tokenString, err := token.SignedString(secret)

	return tokenString, err

}
