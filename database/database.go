package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/svetlyi/dupsfinder/database/migration"
	"log"
)

var db *sql.DB = nil

func GetDB() *sql.DB {
	return db
}

func CreateDB(dbPath string) {
	var err error

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	migration.RunMigrations(db)
}
