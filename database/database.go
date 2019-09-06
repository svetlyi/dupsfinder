package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/svetlyi/dupsfinder/database/migration"
	"log"
)

var db *sql.DB = nil

func init() {
	createDB()
	migration.RunMigrations(db)
	//defer db.Close()
}

func GetDB() *sql.DB {
	return db
}

func createDB() {
	var err error

	db, err = sql.Open("sqlite3", "./dups.db")
	if err != nil {
		log.Fatal(err)
	}
}
