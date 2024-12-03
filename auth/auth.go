package auth

// Root package is for domain types

// type User struct {
// 	UserResponse
// 	Password string `json:"password" db:"password"`
// }

type User struct {
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type UserResponse struct {
	ID    int64  `json:"id" db:"id"`
	Email string `json:"email" db:"email"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:    u.ID,
		Email: u.Email,
	}
}

type JWT struct {
	ID      int64  `json:"id" db:"id"`
	Access  string `json:"access" db:"access"`
	Refresh string `json:"refresh" db:"refresh"`
}

type AuthService interface {
	User(id int) (*User, error)
	// Users() ([]*User, error)
	// CreateUser(u *User) error
	// DeleteUser(id int) error
}

type AuthController interface {
	User(id int) (*User, error)
	// Users() ([]*User, error)
	// CreateUser(u *User) error
	// DeleteUser(id int) error
}

// type Router interface {
// }

type Error struct {
	Code    int
	Message string
}
