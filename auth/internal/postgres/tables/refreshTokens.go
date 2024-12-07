package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func CreateRefreshTokensTable(db *sql.DB) error {
	// ref: https://stackoverflow.com/questions/59691425/how-to-enforce-that-there-is-only-one-true-value-in-a-column-per-names-in-an
	query := `
        CREATE TABLE IF NOT EXISTS refresh_tokens (
            uuid UUID PRIMARY KEY,
            hashed_token TEXT NOT NULL,
            user_uuid UUID NOT NULL,
            active BOOLEAN NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE NOT NULL,
            FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
        );
        
        -- Creates a partial index which ensures there is only a single 
        CREATE UNIQUE INDEX IF NOT EXISTS idx_single_active_token_per_user
        ON refresh_tokens (user_uuid)
        WHERE active = true;`

	_, err := db.Exec(query)
	return err
}
