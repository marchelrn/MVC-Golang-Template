package migrations

import (
	"database/sql"
	"log"
)

type createRefreshTokensTable struct{}

func (m *createRefreshTokensTable) SkipProd() bool {
	return false
}

func getCreateRefreshTokensTable() migration {
	return &createRefreshTokensTable{}
}

func (m *createRefreshTokensTable) Name() string {
	return "create-refresh-tokens"
}

func (m *createRefreshTokensTable) Up(conn *sql.Tx) error {
	log.Println("Creating Up migration: create-refresh-tokens")
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS public.refresh_tokens (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			token_hash TEXT NOT NULL UNIQUE,
			device_id VARCHAR(128) NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			revoked_at TIMESTAMP NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_device_id ON refresh_tokens(device_id);
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_revoked_at ON refresh_tokens(revoked_at);
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_device ON refresh_tokens(user_id, device_id);
	`)
	return err
}

func (m *createRefreshTokensTable) Down(conn *sql.Tx) error {
	log.Println("Creating Down migration: create-refresh-tokens")
	_, err := conn.Exec(`DROP TABLE IF EXISTS refresh_tokens`)
	return err
}
