package auth

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

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
