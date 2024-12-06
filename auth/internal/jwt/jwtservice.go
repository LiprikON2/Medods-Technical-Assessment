package jwt

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"

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

func (j *JWTService) GenerateTokens(payload auth.JWTPayloadDto, accessExpireIn time.Duration, refreshExpireIn time.Duration) (accessToken string, refreshToken string, err error) {
	issuedAt := time.Unix(payload.Iat, 0)

	accessExpireTime := issuedAt.Add(accessExpireIn).Unix()
	accessToken, err = j.newAccessToken(payload, accessExpireTime)
	if err != nil {
		return "", "", err
	}

	refreshExpireTime := issuedAt.Add(refreshExpireIn).Unix()
	refreshToken, err = j.newRefreshToken(payload, refreshExpireTime)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (j *JWTService) newAccessToken(payload auth.JWTPayloadDto, expireTime int64) (string, error) {
	return j.newToken(payload, expireTime, j.accessSecret)
}

func (j *JWTService) newRefreshToken(payload auth.JWTPayloadDto, expireTime int64) (string, error) {
	return j.newToken(payload, expireTime, j.refreshSecret)
}

func (j *JWTService) newToken(payload auth.JWTPayloadDto, expireTime int64, secret []byte) (string, error) {
	mapClaims := jwt.MapClaims{
		"ip":  payload.IP,
		"iat": payload.Iat,
		"exp": expireTime,
		// "nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, mapClaims)
	tokenString, err := token.SignedString(secret)

	return tokenString, err

}
