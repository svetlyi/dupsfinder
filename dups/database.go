package dups

import (
	"database/sql"
	"log"
)

func getInsertStmt(db *sql.DB) *sql.Stmt {
	insertStmt, insertErr := db.Prepare(`INSERT INTO files('path', 'hash') VALUES (?, ?)`)
	if insertErr != nil {
		log.Fatalf("Error in getInsertStmt: %s\n", insertErr)
	}
	return insertStmt
}
