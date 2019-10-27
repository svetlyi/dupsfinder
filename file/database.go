package file

import (
	"database/sql"
	"log"
)

func GetSelectHashByPathStmt(db *sql.DB) *sql.Stmt {
	stmt, err := db.Prepare("SELECT hash FROM files WHERE path=?")
	if nil != err {
		log.Fatalf("Error in GetSelectHashByPathStmt: %s\n", err)
	}

	return stmt
}

func GetHashByPathFromDB(stmt *sql.Stmt, path string) (string, error) {
	var hash string
	selectErr := stmt.QueryRow(path).Scan(&hash)
	if selectErr != nil {
		return "", selectErr
	}

	return hash, nil
}

// GetInsertStmt inserts information about
// a file (path, hash and size)
func GetInsertStmt(db *sql.DB) *sql.Stmt {
	insertStmt, insertErr := db.Prepare(`INSERT INTO files('path', 'hash', 'size') VALUES (?, ?, ?)`)
	if insertErr != nil {
		log.Fatalf("Error in GetInsertStmt: %s\n", insertErr)
	}
	return insertStmt
}

// GetUpdateHashStmt updates hash by path
func GetUpdateHashStmt(db *sql.DB) *sql.Stmt {
	updateStmt, insertErr := db.Prepare(`UPDATE files SET 'hash' = ? WHERE 'path' = ?`)
	if insertErr != nil {
		log.Fatalf("Error in GetUpdateHashStmt: %s\n", insertErr)
	}
	return updateStmt
}
