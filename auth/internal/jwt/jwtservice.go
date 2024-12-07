package jwt

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/netip"

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

	refreshToken = j.newRefreshToken(refreshPayload)

	return accessToken, refreshToken, nil
}

func (j *JWTService) newAccessToken(payload *auth.AccessPayload) (string, error) {
	mapClaims := jwt.MapClaims{
		"jti": payload.Jti,
		"ip":  payload.IP,
		"iat": payload.Iat,
		"sub": payload.Sub,
		"exp": payload.Exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, mapClaims)
	tokenString, err := token.SignedString(j.accessSecret)

	return tokenString, err

}

func (j *JWTService) newRefreshToken(payload *auth.RefreshPayload) string {
	payloadJti := payload.Jti[:]
	payloadIp := payload.IP.AsSlice()
	payloadBytes := append(payloadJti, payloadIp...)

	refreshToken := base64.StdEncoding.EncodeToString(payloadBytes)
	return refreshToken
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

	return j.parseAccessTokenClaims(claims)
}

func (j *JWTService) parseAccessTokenClaims(claims jwt.MapClaims) (*auth.AccessPayload, error) {
	payload := &auth.AccessPayload{}

	if jtiStr, ok := claims["jti"].(string); ok {
		jti, err := j.uuidService.Parse(jtiStr)
		if err != nil {
			return nil, fmt.Errorf("invalid jti claim type")
		}
		payload.Jti = jti
	} else {
		return nil, fmt.Errorf("invalid jti claim type")
	}

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

	return j.parseRefreshToken(paySignCombined)
}

func (j *JWTService) parseRefreshToken(payloadBytes []byte) (*auth.RefreshPayload, error) {
	payloadJti := payloadBytes[:16]

	jti, err := j.uuidService.FromBytes(payloadJti)
	if err != nil {
		return nil, fmt.Errorf("invalid jti type: %w", err)
	}

	payloadIp := payloadBytes[16:]

	ip, ok := netip.AddrFromSlice(payloadIp)
	if !ok && len(payloadIp) != 0 {
		return nil, fmt.Errorf("invalid ip type")
	}

	payload := &auth.RefreshPayload{Jti: jti, IP: ip}
	return payload, nil
}
