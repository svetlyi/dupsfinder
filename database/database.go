package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/svetlyi/dupsfinder/database/migration"
	"log"
)

var db *sql.DB = nil

func NewDB(dbPath string) *sql.DB {
	var err error

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	migration.RunMigrations(db)

	return db
}
