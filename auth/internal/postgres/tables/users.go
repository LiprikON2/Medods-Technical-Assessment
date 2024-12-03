package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func CreateUsersTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS users (
            id BIGSERIAL PRIMARY KEY,
            email TEXT NOT NULL,
            password TEXT NOT NULL,
            CONSTRAINT users_email_unique UNIQUE (email)
        );
        
        CREATE INDEX IF NOT EXISTS users_email_idx ON users(email);
    `

	_, err := db.Exec(query)
	return err
}
