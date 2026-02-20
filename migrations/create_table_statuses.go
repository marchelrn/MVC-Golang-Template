package migrations

import (
	"database/sql"
	"log"
)

type createStatusesTable struct{}

func (m *createStatusesTable) SkipProd() bool {
	return false
}

func getCreateStatusesTable() migration {
	return &createStatusesTable{}
}

func (m *createStatusesTable) Name() string {
	return "create-statuses"
}

func (m *createStatusesTable) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS public.statuses (
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR(50) UNIQUE NOT NULL,
			description VARCHAR(255),
			category VARCHAR(50) NOT NULL,
			color VARCHAR(20),
			order_index INT NOT NULL
		);
	`)
	log.Println("Creating Up migration: create-statuses")
	return err
}

func (m *createStatusesTable) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`DROP TABLE IF EXISTS statuses`)
	if err != nil {
		return err
	}
	return err
}

