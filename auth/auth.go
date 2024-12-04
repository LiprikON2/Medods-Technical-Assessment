package auth

// Root package with domain types
type User struct {
	UserDto
	Password string `json:"password" db:"password"`
}

type UserDto struct {
	ID    int64  `json:"id" db:"id"`
	Email string `json:"email" db:"email"`
}

func (u *User) ToDto() UserDto {
	return UserDto{
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
	Users() ([]*User, error)
	// CreateUser(u *User) error
	// DeleteUser(id int) error
}

type AuthController interface {
	User(id int) (*User, error)
	Users() ([]*User, error)
	// CreateUser(u *User) error
	// DeleteUser(id int) error
}

// type Router interface {
// }

type Error struct {
	Code    int
	Message string
}
