package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
	auth "github.com/medods-technical-assessment"
	"github.com/medods-technical-assessment/internal/common"
	tables "github.com/medods-technical-assessment/internal/postgres/tables"
)

// AuthService represents a PostgreSQL implementation of auth.AuthService.
type AuthService struct {
	DB *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{
		DB: db,
	}
}

func Open(host, port, dbname, user, password string) (*sql.DB, error) {

	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		user,
		dbname,
		password,
	)
	log.Print("Connecting to postgres database...\n", conn)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Panic(err)
	}
	// Create tables
	if err := tables.CreateUsersTable(db); err != nil {
		log.Panic(err)
	}
	if err := tables.CreateRefreshTokensTable(db); err != nil {
		log.Panic(err)
	}

	return db, err

}

func (s *AuthService) GetUser(uuid auth.UUID) (*auth.User, error) {
	user := &auth.User{}
	query := `
        SELECT uuid, email, password
        FROM users
        WHERE uuid = $1`

	err := s.DB.QueryRow(query, uuid).Scan(
		&user.UUID,
		&user.Email,
		&user.Password,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}
	return user, err
}

func (s *AuthService) GetUserByEmail(email string) (*auth.User, error) {
	user := &auth.User{}
	query := `
        SELECT uuid, email, password
        FROM users
        WHERE email = $1`

	err := s.DB.QueryRow(query, email).Scan(
		&user.UUID,
		&user.Email,
		&user.Password,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user with email %v not found: %w", email, err)
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching user with email %v: %w", email, err)
	}
	return user, err
}

func (s *AuthService) GetUsers() ([]*auth.User, error) {
	users := make([]*auth.User, 0)
	query := `
        SELECT uuid, email, password 
        FROM users`

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := &auth.User{}
		err := rows.Scan(
			&user.UUID,
			&user.Email,
			&user.Password,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

func (s *AuthService) CreateUser(user *auth.User) (*auth.User, error) {
	query := `
        INSERT INTO users (uuid, email, password)
        VALUES ($1, $2, $3)
        RETURNING uuid, email, password`

	err := s.DB.QueryRow(
		query,
		user.UUID,
		user.Email,
		user.Password,
	).Scan(
		&user.UUID,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case PgErrUniqueViolation:
				if pqErr.Constraint == common.ConstraintUserEmailUnique {
					return nil, common.ErrDuplicateEmail
				}
			}
		}

		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
}

func (s *AuthService) UpdateUser(user *auth.User) (*auth.User, error) {
	query := `
        UPDATE users
		SET email = $2,
			password = $3
		WHERE uuid = $1
		RETURNING uuid, email, password`

	err := s.DB.QueryRow(
		query,
		user.UUID,
		user.Email,
		user.Password,
	).Scan(
		&user.UUID,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case PgErrUniqueViolation:
				if pqErr.Constraint == common.ConstraintUserEmailUnique {
					return nil, common.ErrDuplicateEmail
				}
			}
		}

		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return user, nil
}

func (s *AuthService) DeleteUser(uuid auth.UUID) error {
	query := `
        DELETE FROM users
		WHERE uuid = $1`

	result, err := s.DB.Exec(query, uuid)

	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("error deleting user: user not found")
	}
	return err
}

func (s *AuthService) AddRefreshToken(refreshToken *auth.RefreshToken) error {
	query := `
        INSERT INTO refresh_tokens (uuid, hashed_token, user_uuid, active, created_at)
        VALUES ($1, $2, $3, $4, $5)`

	_, err := s.DB.Exec(
		query,
		refreshToken.UUID,
		refreshToken.HashedToken,
		refreshToken.UserUUID,
		refreshToken.Active,
		refreshToken.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("error adding refresh token: %w", err)
	}

	return nil
}

func (s *AuthService) RevokeRefreshTokensByUser(userUUID auth.UUID) error {
	query := `
		UPDATE refresh_tokens
		SET active = false
		WHERE user_uuid = $1`

	_, err := s.DB.Exec(query, userUUID)

	if err != nil {
		return fmt.Errorf("error while revoking refresh tokens for user: %w", err)
	}

	return nil
}
func (s *AuthService) GetActiveRefreshTokenByUser(userUUID auth.UUID) (*auth.RefreshToken, error) {
	refreshToken := &auth.RefreshToken{}
	query := `
        SELECT uuid, hashed_token, user_uuid, active, created_at
        FROM refresh_tokens
        WHERE user_uuid = $1 AND
			  active = true`

	err := s.DB.QueryRow(query, userUUID).Scan(
		&refreshToken.UUID,
		&refreshToken.HashedToken,
		&refreshToken.UserUUID,
		&refreshToken.Active,
		&refreshToken.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("refresh token not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("error getting refresh token: %w", err)
	}

	return refreshToken, nil
}

func (s *AuthService) GetActiveRefreshToken(uuid auth.UUID) (*auth.RefreshToken, error) {
	refreshToken := &auth.RefreshToken{}
	query := `
        SELECT uuid, hashed_token, user_uuid, active, created_at
        FROM refresh_tokens
        WHERE uuid = $1 AND
			  active = true`

	err := s.DB.QueryRow(query, uuid).Scan(
		&refreshToken.UUID,
		&refreshToken.HashedToken,
		&refreshToken.UserUUID,
		&refreshToken.Active,
		&refreshToken.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("refresh token not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("error getting refresh token: %w", err)
	}

	return refreshToken, nil
}
