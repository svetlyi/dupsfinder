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
	var filesWithChangedHashes = make(map[string][2]string)
	selectStmt := file.GetSelectHashByPathStmt(db)
	defer selectStmt.Close()
	insertStmt := file.GetInsertStmt(db)
	defer insertStmt.Close()

	for fileInfo := range *filesInfoChannel {
		hashInDB, hashInDBErr := file.GetHashByPathFromDB(selectStmt, fileInfo.Path)
		if nil == hashInDBErr {
			log.Printf("File %s already exists in DB\n", fileInfo.Path)
			if hashInDB != fileInfo.Hash {
				log.Printf("File %s was changed\n", fileInfo.Path)
				filesWithChangedHashes[fileInfo.Path] = [2]string{fileInfo.Hash, hashInDB}
			}
		} else if hashInDBErr == sql.ErrNoRows {
			_, err := insertStmt.Exec(fileInfo.Path, fileInfo.Hash)
			if nil != err {
				log.Fatal(err)
			}
		} else {
			log.Fatal(hashInDBErr)
		}
		if nil == fileDups[fileInfo.Hash] {
			fileDups[fileInfo.Hash] = make([]structs.FileInfo, 0)
		}
		fileDups[fileInfo.Hash] = append(fileDups[fileInfo.Hash], fileInfo)
	}

	printDups(fileDups)
	printChangedFiles(filesWithChangedHashes)
	*doneChannel <- true
}

func printChangedFiles(files map[string][2]string) {
	fmt.Println("==================================")
	fmt.Println("Changed files:")

	for path, hashes := range files {
		fmt.Printf("File %s has hashes:\n", path)

		for _, hash := range hashes {
			fmt.Println(hash)
		}
	}
}

func printDups(files map[string][]structs.FileInfo) {
	fmt.Println("==================================")
	fmt.Println("Dups:")

	for _, sameHashFiles := range files {
		if len(sameHashFiles) > 1 {
			fmt.Printf("Found dups for %v: \n", sameHashFiles[0].Path)

			for _, fileInfo := range sameHashFiles {
				fmt.Println(fileInfo.Path)
			}
		}
	}
}
