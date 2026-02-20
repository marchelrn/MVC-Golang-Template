package migrations

import (
	"database/sql"
	"log"
)

type createEpicsTable struct{}

func (m *createEpicsTable) SkipProd() bool {
	return false
}

func getCreateEpicsTable() migration {
	return &createEpicsTable{}
}

func (m *createEpicsTable) Name() string {
	return "create-epics"
}

func (m *createEpicsTable) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS public.epics (
			id BIGSERIAL PRIMARY KEY,
			project_id BIGINT NOT NULL REFERENCES projects(id),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			status VARCHAR(50) NOT NULL DEFAULT 'open',
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);
	`)
	log.Println("Creating Up migration: create-epics")
	return err
}

func (m *createEpicsTable) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`DROP TABLE IF EXISTS epics`)
	if err != nil {
		return err
	}
	return err
}

