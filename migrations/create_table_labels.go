package migrations

import (
	"database/sql"
	"log"
)

type createLabelsTable struct{}

func (m *createLabelsTable) SkipProd() bool {
	return false
}

func getCreateLabelsTable() migration {
	return &createLabelsTable{}
}

func (m *createLabelsTable) Name() string {
	return "create-labels"
}

func (m *createLabelsTable) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS public.labels (
			id BIGSERIAL PRIMARY KEY,
			project_id BIGINT NOT NULL REFERENCES projects(id),
			name VARCHAR(100) NOT NULL,
			color VARCHAR(20),
			UNIQUE (project_id, name)
		);
	`)
	log.Println("Creating Up migration: create-labels")
	return err
}

func (m *createLabelsTable) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`DROP TABLE IF EXISTS labels`)
	if err != nil {
		return err
	}
	return err
}

