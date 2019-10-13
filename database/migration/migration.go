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
	hash VARCHAR(255),
	size INTEGER
)
`)
}

func RunMigrations(db *sql.DB) {
	for _, query := range queries {
		if _, err := (*db).Exec(query); err != nil {
			log.Fatalf("%q: %s\n", err, query)
		}
	}
	log.Println("Migrated successfully")
}
