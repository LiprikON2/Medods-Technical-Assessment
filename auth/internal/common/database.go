package common

import "fmt"

var (
	ErrDuplicateEmail = fmt.Errorf("user with this email already exists")
)

const (
	ConstraintUserEmailUnique = "users_email_unique"
)
