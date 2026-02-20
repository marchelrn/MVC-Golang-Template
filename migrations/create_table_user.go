package migrations

import (
	"database/sql"
	"log"
)

type createUserTable struct{}

func (m *createUserTable) SkipProd() bool {
	return false
}

func getCreateUserTable() migration {
	return &createUserTable{}
}

func (m *createUserTable) Name() string {
	return "create-user"
}

func (m *createUserTable) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS public.users (
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			username VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'member',
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW(),
			deleted_at TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
	`)
	log.Println("Creating Up migration: create-user")
	return err
}

func (m *createUserTable) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`DROP TABLE IF EXISTS users`)
	if err != nil {
		return err
	}
	return err
}
