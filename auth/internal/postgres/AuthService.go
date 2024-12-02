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

// User returns a user for a given id.
func (s *AuthService) User(id int) (*auth.User, error) {
	var u auth.User
	// row := db.QueryRow(`SELECT id, email FROM users WHERE id = $1`, id)
	// if row.Scan(&u.ID, &u.Email); err != nil {
	// 	return nil, err
	// }
	return &u, nil
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
	log.Print(conn)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	// db.AutoMigrate(&domain.Message{})

	return db, err

}
