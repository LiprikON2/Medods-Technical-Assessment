package auth

// Root package with domain types

import (
	"net/http"
	"net/netip"
	"time"

	"github.com/google/uuid"
)

type UUID = uuid.UUID

type User struct {
	UUID          UUID            `json:"uuid" db:"uuid"`
	Email         string          `json:"email" db:"email"`
	Password      string          `json:"-" db:"password"`
	RefreshTokens []*RefreshToken `json:"-" db:"refresh_tokens"`
}

type PublicUser struct {
	Email string `json:"email"`
}

func (u User) ToPublic() *PublicUser {
	return &PublicUser{
		Email: u.Email,
	}
}

type CreateUserDto struct {
	Email    string `json:"email" validate:"required,email,max=254"`
	Password string `json:"password" validate:"required,password,min=8"`
}

type UpdateUserDto struct {
	Email    string `json:"email" validate:"omitempty,email,max=254"`
	Password string `json:"password" validate:"omitempty,password,min=8"`
}

type LoginUserDto struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type RefreshToken struct {
	UUID        UUID      `json:"uuid" db:"uuid"`
	HashedToken string    `json:"hashedToken" db:"hashed_token"`
	UserUUID    UUID      `json:"userUUID" db:"user_uuid"`
	Active      bool      `json:"active" db:"active"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AuthService interface {
	GetUser(uuid UUID) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUsers() ([]*User, error)
	CreateUser(user *User) (*User, error)
	UpdateUser(user *User) (*User, error)
	DeleteUser(uuid UUID) error
	AddRefreshToken(refreshToken *RefreshToken) error
	RevokeRefreshTokensByUser(userUUID UUID) error
	GetActiveRefreshTokenByUser(userUUID UUID) (*RefreshToken, error)
	GetActiveRefreshToken(uuid UUID) (*RefreshToken, error)
}

type AuthController interface {
	GetUser(w http.ResponseWriter, r *http.Request)
	GetMe(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	LoginByUUID(w http.ResponseWriter, r *http.Request)
	GetUsers(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type ValidationService interface {
	ValidateUserInput(input any) []ValidationError
}

type RequestError struct {
	Code    int `json:"code"`
	Message any `json:"message"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type CryptoService interface {
	HashPassword(password string) string
	ComparePasswords(hpass string, pass string) error
}

type UUIDService interface {
	New() UUID
	Parse(s string) (UUID, error)
	MustParse(s string) UUID
	FromBytes(b []byte) (uuid UUID, err error)
}

type JWTService interface {
	GenerateTokens(refreshPayload *RefreshPayload, accessPayload *AccessPayload) (accessToken string, refreshToken string, err error)
	VerifyAccessToken(accessToken string) error
	GetAccessTokenPayload(accessToken string) (*AccessPayload, error)
	GetRefreshTokenPayload(refreshToken string) (*RefreshPayload, error)
}
type RefreshPayload struct {
	// Unique identifier which ensures that:
	// - each generated token is unique
	// - refresh and access tokens are coupled
	Jti UUID `json:"jti"`
	// User's IPv4 or IPv6 address (without port)
	IP netip.Addr `json:"ip"`
}
type AccessPayload struct {
	// Unique identifier which ensures that:
	// - each generated token is unique
	// - refresh and access tokens are coupled
	Jti UUID `json:"jti"`
	// User's IPv4 or IPv6 address (without port)
	IP string `json:"ip"`
	// Issued at
	Iat int64 `json:"iat"`
	// Expiration time
	Exp int64 `json:"exp"`
}

type MailService interface {
	Send(to, subject, message string) error
}
