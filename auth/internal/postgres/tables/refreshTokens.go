package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func CreateRefreshTokensTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS refresh_tokens (
            uuid UUID PRIMARY KEY,
            hashed_token TEXT NOT NULL,
            user_uuid UUID NOT NULL,
            revoked BOOLEAN NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE NOT NULL,
            FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
        );`

	_, err := db.Exec(query)
	return err
}
