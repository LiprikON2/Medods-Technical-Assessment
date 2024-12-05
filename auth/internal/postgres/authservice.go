package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// Autoloads `.env`
	_ "github.com/joho/godotenv/autoload"
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

func (s *AuthService) User(id int) (*auth.User, error) {
	user := &auth.User{}
	query := `
        SELECT id, email, password
        FROM users
        WHERE id = $1`

	err := s.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user with id %v not found: %w", id, err)
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching user with id %v: %w", id, err)
	}
	return user, err
}

func (s *AuthService) Users() ([]*auth.User, error) {
	users := make([]*auth.User, 0)
	query := `
        SELECT id, email, password 
        FROM users`

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := &auth.User{}
		err := rows.Scan(
			&user.ID,
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
        INSERT INTO users (email, password)
        VALUES ($1, $2)
        RETURNING id, email, password`

	err := s.DB.QueryRow(
		query,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
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
		WHERE id = $1
		RETURNING id, email, password`

	err := s.DB.QueryRow(
		query,
		user.ID,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
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

func Open() (*sql.DB, error) {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DATABASE")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

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
		log.Fatal(err)
	}

	return db, err

}
