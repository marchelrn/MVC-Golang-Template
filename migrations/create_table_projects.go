package migrations

import (
	"database/sql"
	"log"
)

type createProjectsTable struct{}

func (m *createProjectsTable) SkipProd() bool {
	return false
}

func getCreateProjectsTable() migration {
	return &createProjectsTable{}
}

func (m *createProjectsTable) Name() string {
	return "create-projects"
}

func (m *createProjectsTable) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS public.projects (
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			key VARCHAR(10) UNIQUE NOT NULL,
			description TEXT,
			lead_id BIGINT REFERENCES users(id),
			created_by BIGINT REFERENCES users(id),
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW(),
			deleted_at TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_projects_key ON projects(key);
		CREATE INDEX IF NOT EXISTS idx_projects_lead ON projects(lead_id);
	`)
	log.Println("Creating Up migration: create-projects")
	return err
}

func (m *createProjectsTable) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`DROP TABLE IF EXISTS projects`)
	if err != nil {
		return err
	}
	return err
}

