package dups

import (
	"database/sql"
	"fmt"
	"github.com/svetlyi/dupsfinder/file"
	"github.com/svetlyi/dupsfinder/structs"
	"log"
)

func ListenFilesInfoChannel(filesInfoChannel *chan structs.FileInfo, doneChannel *chan bool, db *sql.DB) {
	var fileDups = make(map[string][]structs.FileInfo)
	insertStmt := file.GetInsertStmt(db)
	defer insertStmt.Close()

	for fileInfo := range *filesInfoChannel {
		if nil == fileDups[fileInfo.Hash] {
			fileDups[fileInfo.Hash] = make([]structs.FileInfo, 0)
		}
		_, err := insertStmt.Exec(fileInfo.Path, fileInfo.Hash)
		if nil != err {
			log.Fatal(err)
		}
		fileDups[fileInfo.Hash] = append(fileDups[fileInfo.Hash], fileInfo)
	}

	printDups(fileDups)
	*doneChannel <- true
}

func printDups(files map[string][]structs.FileInfo) {
	fmt.Println("==================================")
	fmt.Println("Dups:")

	for _, sameHashFiles := range files {
		if len(sameHashFiles) > 1 {
			fmt.Printf("Found dups for %v: \n", sameHashFiles[0].Path)

			for _, file := range sameHashFiles {
				fmt.Println(file.Path)
			}
		}
	}
}
