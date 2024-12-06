package jwt

import (
	"encoding/base64"
	"fmt"
	"log"

	auth "github.com/medods-technical-assessment"

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
	accessSecret, err := base64.StdEncoding.DecodeString(accessSecretStr)
	if err != nil {
		log.Panic(fmt.Errorf("couldn't convert accessSecret into bytes: %w", err))
	}

	return &JWTService{
		refreshSecret: refreshSecret,
		accessSecret:  accessSecret,
	}
}

func (j *JWTService) NewAccessToken(payload auth.JWTPayload) (string, error) {
	return j.newToken(payload, j.accessSecret)
}

func (j *JWTService) NewRefreshToken(payload auth.JWTPayload) (string, error) {
	return j.newToken(payload, j.refreshSecret)
}

func (j *JWTService) GenerateTokens(payload auth.JWTPayload) (accessToken string, refreshToken string, err error) {
	accessToken, err = j.NewAccessToken(payload)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = j.NewRefreshToken(payload)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (j *JWTService) newToken(payload auth.JWTPayload, secret []byte) (string, error) {
	mapClaims := jwt.MapClaims{
		"ip":  payload.IP,
		"iat": payload.Iat,
		"exp": payload.Exp,
		// "nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, mapClaims)
	tokenString, err := token.SignedString(secret)

	return tokenString, err

}
