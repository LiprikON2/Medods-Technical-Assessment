package jwt

import (
	"encoding/base64"
	"fmt"
	"log"

	auth "github.com/medods-technical-assessment"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	uuidService  auth.UUIDService
	accessSecret []byte
}

func NewJWTService(accessSecretStr string, uuidService auth.UUIDService) *JWTService {

	if accessSecretStr == "" {
		log.Panic(fmt.Errorf("accessSecretStr is empty"))
	}

	// ref: https://golang-jwt.github.io/jwt/usage/signing_methods/#signing-methods-and-key-types
	accessSecret, err := base64.StdEncoding.DecodeString(accessSecretStr)
	if err != nil {
		log.Panic(fmt.Errorf("couldn't convert accessSecret into bytes: %w", err))
	}

	return &JWTService{
		accessSecret: accessSecret,
		uuidService:  uuidService,
	}
}

func (j *JWTService) GenerateTokens(refreshPayload *auth.RefreshPayload, accessPayload *auth.AccessPayload) (accessToken string, refreshToken string, err error) {
	accessToken, err = j.newAccessToken(accessPayload)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = j.newRefreshToken(refreshPayload)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (j *JWTService) newAccessToken(payload *auth.AccessPayload) (string, error) {
	mapClaims := jwt.MapClaims{
		"ip":  payload.IP,
		"iat": payload.Iat,
		"sub": payload.Sub,
		"exp": payload.Exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, mapClaims)
	tokenString, err := token.SignedString(j.accessSecret)

	return tokenString, err

}

// Returns byte array slice of payload + HS256 signature encoded in base64 string
// TODO since token is hashed with bcrypt, signature might not be needed
func (j *JWTService) newRefreshToken(payload *auth.RefreshPayload) (string, error) {
	payloadJti := payload.Jti[:]

	// payloadBytes := append(payloadJti, signature...)
	payloadBytes := payloadJti

	refreshToken := base64.StdEncoding.EncodeToString(payloadBytes)
	return refreshToken, nil
}

func (j *JWTService) VerifyAccessToken(accessToken string) error {
	_, err := j.getAccessTokenPayload(accessToken, j.accessSecret)
	return err
}

func (j *JWTService) GetAccessTokenPayload(accessToken string) (*auth.AccessPayload, error) {
	return j.getAccessTokenPayload(accessToken, j.accessSecret)
}

func (j *JWTService) getAccessTokenPayload(tokenString string, secret []byte) (*auth.AccessPayload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	payload, err := j.parseAccessTokenClaims(claims)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (j *JWTService) parseAccessTokenClaims(claims jwt.MapClaims) (*auth.AccessPayload, error) {
	payload := &auth.AccessPayload{}

	if ip, ok := claims["ip"].(string); ok {
		payload.IP = ip
	} else {
		return nil, fmt.Errorf("invalid ip claim type")
	}

	if exp, ok := claims["exp"].(float64); ok {
		payload.Exp = int64(exp)
	} else {
		return nil, fmt.Errorf("invalid exp claim type")
	}

	if iat, ok := claims["iat"].(float64); ok {
		payload.Iat = int64(iat)
	} else {
		return nil, fmt.Errorf("invalid iat claim type")
	}

	if sub, ok := claims["sub"].(string); ok {
		userUUID, err := j.uuidService.Parse(sub)
		if err != nil {
			return nil, fmt.Errorf("invalid sub claim type")
		}
		payload.Sub = userUUID

	} else {
		return nil, fmt.Errorf("invalid sub claim type")
	}

	return payload, nil
}

func (j *JWTService) GetRefreshTokenPayload(refreshToken string) (*auth.RefreshPayload, error) {
	return j.getRefreshTokenPayload(refreshToken)
}

func (j *JWTService) getRefreshTokenPayload(tokenString string) (*auth.RefreshPayload, error) {
	paySignCombined, err := base64.StdEncoding.DecodeString(tokenString)
	if err != nil {
		return nil, fmt.Errorf("error decoding refresh token: %w", err)
	}

	payload, err := j.parseRefreshToken(paySignCombined)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (j *JWTService) parseRefreshToken(paySignCombined []byte) (*auth.RefreshPayload, error) {
	payload := &auth.RefreshPayload{}
	payloadJti := paySignCombined[:16]
	// signature := paySignCombined[len(paySignCombined)-32:]

	jti, err := j.uuidService.FromBytes(payloadJti)
	if err != nil {
		return nil, fmt.Errorf("invalid jti type: %w", err)
	}

	payload.Jti = jti
	return payload, nil
}
