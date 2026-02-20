package migrations

import (
	"database/sql"
	"log"
)

type createPrioritiesTable struct{}

func (m *createPrioritiesTable) SkipProd() bool {
	return false
}

func getCreatePrioritiesTable() migration {
	return &createPrioritiesTable{}
}

func (m *createPrioritiesTable) Name() string {
	return "create-priorities"
}

func (m *createPrioritiesTable) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS public.priorities (
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR(50) UNIQUE NOT NULL,
			description VARCHAR(255),
			color VARCHAR(20),
			order_index INT NOT NULL
		);
	`)
	log.Println("Creating Up migration: create-priorities")
	return err
}

func (m *createPrioritiesTable) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`DROP TABLE IF EXISTS priorities`)
	if err != nil {
		return err
	}
	return err
}

