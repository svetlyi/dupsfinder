package migration

import (
	"database/sql"
	"log"
)

var queries []string

func init() {
	queries = append(queries, `
CREATE TABLE IF NOT EXISTS files (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	path VARCHAR(255),
	hash VARCHAR(255)
)
`)
}

func RunMigrations(db *sql.DB) {
	for _, query := range queries {
		_, err := (*db).Exec(query)
		if err != nil {
			log.Fatalf("%q: %s\n", err, query)
		}
	}
	log.Println("Migrated successfully")
}
