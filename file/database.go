package file

import (
	"database/sql"
	"log"
)

func GetSelectHashByPathStmt(db *sql.DB) *sql.Stmt {
	selectStmt, selectErr := db.Prepare("SELECT hash FROM files WHERE path=?")

	if nil != selectErr {
		log.Fatalf("Error in GetSelectHashByPathStmt: %s\n", selectErr)
	}
	return selectStmt
}

func GetHashByPathFromDB(stmt *sql.Stmt, path string) (string, error) {
	var hash string
	selectErr := stmt.QueryRow(path).Scan(&hash)
	if selectErr != nil {
		return "", selectErr
	}

	return hash, nil
}

func GetSelectFilesByDir(db *sql.DB) *sql.Stmt {
	selectStmt, selectErr := db.Prepare("SELECT path, hash FROM files WHERE path LIKE ?")

	if nil != selectErr {
		log.Fatalf("Error in GetSelectFilesByDir: %s\n", selectErr)
	}
	return selectStmt
}

func GetInsertStmt(db *sql.DB) *sql.Stmt {
	insertStmt, insertErr := db.Prepare(`INSERT INTO files('path', 'hash') VALUES (?, ?)`)
	if insertErr != nil {
		log.Fatalf("Error in getInsertStmt: %s\n", insertErr)
	}
	return insertStmt
}
