package migrations

import (
	"database/sql"
)

type alterUserTableAddAvatarUrl struct{}

func (m *alterUserTableAddAvatarUrl) SkipProd() bool {
	return false
}

func getAlterUserTableAddAvatarUrl() migration {
	return &alterUserTableAddAvatarUrl{}
}

func (m *alterUserTableAddAvatarUrl) Name() string {
	return "alter-user-add-avatar-url"
}

func (m *alterUserTableAddAvatarUrl) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		ALTER TABLE users
		ADD COLUMN IF NOT EXISTS avatar_url VARCHAR(255) DEFAULT 'https://ik.imagekit.io/x3x3vd0wg/default_pp.svg';
	`)
	return err
}

func (m *alterUserTableAddAvatarUrl) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`
		ALTER TABLE users
		DROP COLUMN IF EXISTS avatar_url
	`)
	return err
}
