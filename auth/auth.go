package auth

import "net/http"

// Root package with domain types
type User struct {
	ID       int64  `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
}

type CreateUserDto struct {
	Email    string `json:"email" validate:"required,email,max=254"`
	Password string `json:"password" validate:"required,password,min=8"`
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
	User(id int) (*User, error)
	Users() ([]*User, error)
	CreateUser(u *User) (*User, error)
	UpdateUser(u *User) (*User, error)
	DeleteUser(id int) error
}

type AuthController interface {
	User(w http.ResponseWriter, r *http.Request)
	Users(w http.ResponseWriter, r *http.Request)
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
