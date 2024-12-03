package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// Autoloads `.env`
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	auth "github.com/medods-technical-assessment"
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

// User returns a user for a given id.
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

// User returns all users
// func (s *AuthService) Users(id int) (*[]auth.User, error) {
// 	var u auth.User
// 	return &u, nil
// }

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
	// db.AutoMigrate(&domain.Message{})

	return db, err

}
