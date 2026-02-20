package migrations

import (
	"database/sql"
)

type alterUserTableAddStatus struct{}

func (m *alterUserTableAddStatus) SkipProd() bool {
	return false
}

func getAlterUserTableAddStatus() migration {
	return &alterUserTableAddStatus{}
}

func (m *alterUserTableAddStatus) Name() string {
	return "alter-user-add-status"
}

func (m *alterUserTableAddStatus) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		ALTER TABLE users
		ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'active';
	`)
	return err
}

func (m *alterUserTableAddStatus) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`
		ALTER TABLE users
		DROP COLUMN IF EXISTS status
	`)
	return err
}
