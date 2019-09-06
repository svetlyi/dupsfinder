package file

import (
	"database/sql"
	"log"
)

func getSelectHashByPathStmt(db *sql.DB) *sql.Stmt {
	selectStmt, selectErr := db.Prepare("SELECT hash FROM files WHERE path=?")

	if nil != selectErr {
		log.Fatalf("Error in getSelectHashByPathStmt: %s\n", selectErr)
	}
	return selectStmt
}
