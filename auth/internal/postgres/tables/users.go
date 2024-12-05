package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/medods-technical-assessment/internal/common"
)

func CreateUsersTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS users (
            uuid UUID PRIMARY KEY,
            email TEXT NOT NULL,
            password TEXT NOT NULL,
            CONSTRAINT %s UNIQUE (email)
        );`

	_, err := db.Exec(fmt.Sprintf(query, common.ConstraintUserEmailUnique))

	return err
}
