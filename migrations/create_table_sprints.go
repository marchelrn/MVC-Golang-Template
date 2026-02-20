package migrations

import (
	"database/sql"
	"log"
)

type createSprintsTable struct{}

func (m *createSprintsTable) SkipProd() bool {
	return false
}

func getCreateSprintsTable() migration {
	return &createSprintsTable{}
}

func (m *createSprintsTable) Name() string {
	return "create-sprints"
}

func (m *createSprintsTable) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS public.sprints (
			id BIGSERIAL PRIMARY KEY,
			project_id BIGINT NOT NULL REFERENCES projects(id),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			start_date DATE NOT NULL,
			end_date DATE NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'planned',
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);
	`)
	log.Println("Creating Up migration: create-sprints")
	return err
}

func (m *createSprintsTable) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`DROP TABLE IF EXISTS sprints`)
	if err != nil {
		return err
	}
	return err
}

