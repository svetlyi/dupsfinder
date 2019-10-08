package dups

import (
	"database/sql"
	"log"
)

func GetSelectDupsByDir(db *sql.DB) *sql.Stmt {
	query := `
SELECT f1.path, f1.hash
FROM files f1
WHERE f1.hash IN (
	SELECT f2.hash
	FROM files f2
	WHERE f2.path LIKE ? || '%'
	GROUP BY f2.hash
	HAVING COUNT(*) > 1
) ORDER BY f1.hash, LENGTH(f1.path) DESC
LIMIT ? OFFSET ?
`
	selectStmt, selectErr := db.Prepare(query)

	if nil != selectErr {
		log.Fatalf("Error in GetSelectDupsByDir: %s\n", selectErr)
	}
	return selectStmt
}
