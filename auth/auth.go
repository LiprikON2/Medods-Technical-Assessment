package auth

type User struct {
	UserResponse
	Password string `json:"password"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:    u.ID,
		Email: u.Email,
	}
}

// type UserParams struct {
// 	UserID int
// }

type JWT struct {
	ID      int    `json:"id"`
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
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

type Error struct {
	Code    int
	Message string
}
