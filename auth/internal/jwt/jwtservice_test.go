package jwt

import (
	"encoding/base64"
	"fmt"
	"net/netip"
	"strings"
	"testing"
	"time"

	auth "github.com/medods-technical-assessment"
	"github.com/medods-technical-assessment/internal/uuid"
)

func TestJWTServiceGenerateTokens(t *testing.T) {
	us := uuid.NewUUIDService()
	type Input struct {
		userUUIDStr   string
		issuedAt      time.Time
		ipStr         string
		jtiStr        string
		expiredAccess bool
	}
	type Want struct {
		isValidGen               bool
		isValidLen               bool
		isValidAccessHeaderEnc   bool
		isValidAccessHeaderDec   bool
		isValidAccessPayloadEnc  bool
		isValidAccessPayloadDec  bool
		isValidRefreshEnc        bool
		isValidRefreshPayloadEnc bool
		isSameRefreshPayloadIP   bool
		isSameRefreshPayloadJti  bool
		isAccessTokenVerified    bool
	}

	var tests = []struct {
		name  string
		input *Input
		want  *Want
	}{
		{"Valid ipv4",
			&Input{
				userUUIDStr:   "163023da-319c-4a87-8051-6ba52631038a",
				issuedAt:      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC),
				ipStr:         "255.255.255.255",
				jtiStr:        "91a5f9ca-8339-4fde-83f9-d4b6ad1b5131",
				expiredAccess: false,
			}, &Want{
				isValidGen:               true,
				isValidLen:               true,
				isValidAccessHeaderEnc:   true,
				isValidAccessHeaderDec:   true,
				isValidAccessPayloadEnc:  true,
				isValidAccessPayloadDec:  true,
				isValidRefreshEnc:        true,
				isValidRefreshPayloadEnc: true,
				isSameRefreshPayloadIP:   true,
				isSameRefreshPayloadJti:  true,
				isAccessTokenVerified:    true,
			},
		},
		{"Valid ipv6",
			&Input{
				userUUIDStr:   "163023da-319c-4a87-8051-6ba52631038a",
				issuedAt:      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC),
				ipStr:         "2001:db8:3333:4444:5555:6666:7777:8888",
				jtiStr:        "91a5f9ca-8339-4fde-83f9-d4b6ad1b5131",
				expiredAccess: false,
			}, &Want{
				isValidGen:               true,
				isValidLen:               true,
				isValidAccessHeaderEnc:   true,
				isValidAccessHeaderDec:   true,
				isValidAccessPayloadEnc:  true,
				isValidAccessPayloadDec:  true,
				isValidRefreshEnc:        true,
				isValidRefreshPayloadEnc: true,
				isSameRefreshPayloadIP:   true,
				isSameRefreshPayloadJti:  true,
				isAccessTokenVerified:    true,
			},
		},
		{"Missing ip",
			&Input{
				userUUIDStr:   "163023da-319c-4a87-8051-6ba52631038a",
				issuedAt:      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC),
				ipStr:         "",
				jtiStr:        "91a5f9ca-8339-4fde-83f9-d4b6ad1b5131",
				expiredAccess: false,
			}, &Want{
				isValidGen:               true,
				isValidLen:               true,
				isValidAccessHeaderEnc:   true,
				isValidAccessHeaderDec:   true,
				isValidAccessPayloadEnc:  true,
				isValidAccessPayloadDec:  true,
				isValidRefreshEnc:        true,
				isValidRefreshPayloadEnc: true,
				isSameRefreshPayloadIP:   true,
				isSameRefreshPayloadJti:  true,
				isAccessTokenVerified:    true,
			},
		},
		{"Invalid ip",
			&Input{
				userUUIDStr:   "163023da-319c-4a87-8051-6ba52631038a",
				issuedAt:      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC),
				ipStr:         "123",
				jtiStr:        "91a5f9ca-8339-4fde-83f9-d4b6ad1b5131",
				expiredAccess: false,
			}, &Want{
				isValidGen:               true,
				isValidLen:               true,
				isValidAccessHeaderEnc:   true,
				isValidAccessHeaderDec:   true,
				isValidAccessPayloadEnc:  true,
				isValidAccessPayloadDec:  true,
				isValidRefreshEnc:        true,
				isValidRefreshPayloadEnc: true,
				isSameRefreshPayloadIP:   true,
				isSameRefreshPayloadJti:  true,
				isAccessTokenVerified:    true,
			},
		},
		{"Expired access token",
			&Input{
				userUUIDStr:   "163023da-319c-4a87-8051-6ba52631038a",
				issuedAt:      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC),
				ipStr:         "255.255.255.255",
				jtiStr:        "91a5f9ca-8339-4fde-83f9-d4b6ad1b5131",
				expiredAccess: true,
			}, &Want{
				isValidGen:               true,
				isValidLen:               true,
				isValidAccessHeaderEnc:   true,
				isValidAccessHeaderDec:   true,
				isValidAccessPayloadEnc:  true,
				isValidAccessPayloadDec:  true,
				isValidRefreshEnc:        true,
				isValidRefreshPayloadEnc: true,
				isSameRefreshPayloadIP:   true,
				isSameRefreshPayloadJti:  true,
				isAccessTokenVerified:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userUUID := us.MustParse(tt.input.userUUIDStr)
			issuedAt := tt.input.issuedAt
			var accessTokenExpireTime time.Duration
			if tt.input.expiredAccess {
				accessTokenExpireTime = time.Since(issuedAt)
			} else {
				accessTokenExpireTime = time.Since(issuedAt) + 5*time.Minute

			}
			ipStr := tt.input.ipStr
			ip, _ := netip.ParseAddr(ipStr)
			jti := us.MustParse(tt.input.jtiStr)

			refreshPayload := &auth.RefreshPayload{Jti: jti, IP: ip}
			accessPayload := &auth.AccessPayload{Jti: refreshPayload.Jti, IP: ipStr, Sub: userUUID, Iat: issuedAt.Unix(), Exp: issuedAt.Add(accessTokenExpireTime).Unix()}

			js := NewJWTService("MTIzNA==", us)
			accessToken, refreshToken, err := js.GenerateTokens(refreshPayload, accessPayload)

			isValidGen := err == nil
			if isValidGen != tt.want.isValidGen {
				t.Errorf("got isValidGen %v, want valid %v", isValidGen, tt.want.isValidGen)
				return
			}

			parts := strings.Split(accessToken, ".")
			isValidLen := len(parts) == 3

			if isValidLen != tt.want.isValidLen {
				t.Errorf("got isValidLen %v, want valid %v", isValidLen, tt.want.isValidLen)
				return
			}
			header := parts[0]
			payload := parts[1]

			accessHeaderBytes, err := base64.RawStdEncoding.DecodeString(header)
			isValidAccessHeaderEnc := err == nil
			if isValidAccessHeaderEnc != tt.want.isValidAccessHeaderEnc {
				t.Errorf("got isValidAccessHeaderEnc %v, want valid %v", isValidAccessHeaderEnc, tt.want.isValidAccessHeaderEnc)
				return
			}
			accessHeaderDecoded := string(accessHeaderBytes)

			isValidAccessHeaderDec :=
				strings.Contains(accessHeaderDecoded, "\"alg\":\"HS512\"") &&
					strings.Contains(accessHeaderDecoded, "\"typ\":\"JWT\"")
			if isValidAccessHeaderDec != tt.want.isValidAccessHeaderDec {
				t.Errorf("got isValidAccessHeaderDec %v, want valid %v", isValidAccessHeaderDec, tt.want.isValidAccessHeaderDec)
				return
			}

			accessPayloadBytes, err := base64.RawStdEncoding.DecodeString(payload)
			isValidAccessPayloadEnc := err == nil
			if isValidAccessPayloadEnc != tt.want.isValidAccessPayloadEnc {
				t.Errorf("got isValidAccessPayloadEnc %v, want valid %v", isValidAccessPayloadEnc, tt.want.isValidAccessPayloadEnc)
				return

			}
			accessPayloadDecoded := string(accessPayloadBytes)

			isValidAccessPayloadDec := strings.Contains(accessPayloadDecoded, fmt.Sprintf("\"exp\":%v", accessPayload.Exp)) &&
				strings.Contains(accessPayloadDecoded, fmt.Sprintf("\"iat\":%v", accessPayload.Iat)) &&
				strings.Contains(accessPayloadDecoded, fmt.Sprintf("\"ip\":\"%v\"", accessPayload.IP)) &&
				strings.Contains(accessPayloadDecoded, fmt.Sprintf("\"sub\":\"%v\"", accessPayload.Sub)) &&
				strings.Contains(accessPayloadDecoded, fmt.Sprintf("\"jti\":\"%v\"", accessPayload.Jti))
			if isValidAccessPayloadDec != tt.want.isValidAccessPayloadDec {
				t.Errorf("got isValidAccessPayloadDec %v, want valid %v", isValidAccessPayloadDec, tt.want.isValidAccessPayloadDec)
				return
			}

			_, err = base64.StdEncoding.DecodeString(refreshToken)
			isValidRefreshEnc := err == nil
			if isValidRefreshEnc != tt.want.isValidRefreshEnc {
				t.Errorf("got isValidRefreshEnc %v, want valid %v", isValidRefreshEnc, tt.want.isValidRefreshEnc)
				return

			}

			refreshPayloadDecoded, err := js.GetRefreshTokenPayload(refreshToken)
			isValidRefreshPayloadEnc := err == nil
			if isValidRefreshPayloadEnc != tt.want.isValidRefreshPayloadEnc {
				t.Errorf("got isValidRefreshPayloadEnc %v, want valid %v", isValidRefreshPayloadEnc, tt.want.isValidRefreshPayloadEnc)
				return
			}

			isSameRefreshPayloadIP := refreshPayloadDecoded.IP == refreshPayload.IP
			if isSameRefreshPayloadIP != tt.want.isSameRefreshPayloadIP {
				t.Errorf("got isSameRefreshPayloadIP %v, want valid %v", isSameRefreshPayloadIP, tt.want.isSameRefreshPayloadIP)
				return
			}
			isSameRefreshPayloadJti := refreshPayloadDecoded.Jti == refreshPayload.Jti
			if isSameRefreshPayloadJti != tt.want.isSameRefreshPayloadJti {
				t.Errorf("got isSameRefreshPayloadJti %v, want valid %v", isSameRefreshPayloadJti, tt.want.isSameRefreshPayloadJti)
				return
			}

			err = js.VerifyAccessToken(accessToken)
			isAccessTokenVerified := err == nil
			if isAccessTokenVerified != tt.want.isAccessTokenVerified {
				t.Errorf("got isAccessTokenVerified %v, want valid %v", isAccessTokenVerified, tt.want.isAccessTokenVerified)
				return
			}

		})
	}
}
