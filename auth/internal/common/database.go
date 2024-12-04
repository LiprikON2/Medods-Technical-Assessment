package common

import "errors"

var (
	ErrDuplicateEmail = errors.New("user with this email already exists")
)

const (
	ConstraintUserEmailUnique = "users_email_unique"
)
