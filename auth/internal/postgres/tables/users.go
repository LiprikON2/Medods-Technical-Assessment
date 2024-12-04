package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/medods-technical-assessment/internal/common"
)

func CreateUsersTable(db *sql.DB) error {
	query := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS users (
            id BIGSERIAL PRIMARY KEY,
            email TEXT NOT NULL,
            password TEXT NOT NULL,
            CONSTRAINT %s UNIQUE (email)
        );
        
        CREATE INDEX IF NOT EXISTS users_email_idx ON users(email);
    `, common.ConstraintUserEmailUnique)

	_, err := db.Exec(query)
	return err
}
