package auth

import (
	"net/http"

	"github.com/google/uuid"
)

// Root package with domain types
type User struct {
	UUID     UUID   `json:"uuid" db:"uuid"`
	Email    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
}

type CreateUserDto struct {
	Email    string `json:"email" validate:"required,email,max=254"`
	Password string `json:"password" validate:"required,password,min=8"`
}

type Login struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type UpdateUserDto struct {
	Email    string `json:"email" validate:"omitempty,email,max=254"`
	Password string `json:"password" validate:"omitempty,password,min=8"`
}

type JWT struct {
	ID      int64  `json:"id" db:"id"`
	Access  string `json:"access" db:"access"`
	Refresh string `json:"refresh" db:"refresh"`
}

type AuthService interface {
	GetUser(uuid UUID) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUsers() ([]*User, error)
	CreateUser(u *User) (*User, error)
	UpdateUser(u *User) (*User, error)
	DeleteUser(uuid UUID) error
}

type AuthController interface {
	GetUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	GetUsers(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
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
}

// This somewhat breaks DDD
// ref: https://stackoverflow.com/a/31933240
type UUID = uuid.UUID
