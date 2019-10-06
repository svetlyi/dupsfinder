package dups

import (
	"database/sql"
	"fmt"
	"github.com/svetlyi/dupsfinder/file"
	"github.com/svetlyi/dupsfinder/log"
	"github.com/svetlyi/dupsfinder/structs"
)

// storeFilesInfo listens to filesInfoChannel and stores
// the information in a databse
func storeFilesInfo(filesInfoChannel *chan structs.FileInfo, app *structs.App) {
	//file duplicates with the same hash
	var fileDups = make(map[string][]structs.FileInfo)
	//[2]string - two hashes
	var filesWithChangedHashes = make(map[string][2]string)
	selectStmt := file.GetSelectHashByPathStmt(app.DB)
	defer selectStmt.Close()
	insertStmt := file.GetInsertStmt(app.DB)
	defer insertStmt.Close()

	for fileInfo := range *filesInfoChannel {
		hashInDB, hashInDBErr := file.GetHashByPathFromDB(selectStmt, fileInfo.Path)
		if nil == hashInDBErr {
			log.Msg(app.LogChan, fmt.Sprintf("File %s already exists in DB\n", fileInfo.Path))
			if hashInDB != fileInfo.Hash {
				log.Msg(app.LogChan, fmt.Sprintf("File %s was changed\n", fileInfo.Path))
				filesWithChangedHashes[fileInfo.Path] = [2]string{fileInfo.Hash, hashInDB}
			}
		} else if hashInDBErr == sql.ErrNoRows {
			if _, err := insertStmt.Exec(fileInfo.Path, fileInfo.Hash); nil != err {
				log.Err(app.LogChan, err.Error())
				close(*app.ExitChan)
			}
		} else {
			log.Err(app.LogChan, hashInDBErr.Error())
			close(*app.ExitChan)
		}
		if nil == fileDups[fileInfo.Hash] {
			fileDups[fileInfo.Hash] = make([]structs.FileInfo, 0)
		}
		fileDups[fileInfo.Hash] = append(fileDups[fileInfo.Hash], fileInfo)
	}

	printDups(app.LogChan, fileDups)
	printChangedFiles(app.LogChan, filesWithChangedHashes)
	*app.DoneChan <- true
}

func printChangedFiles(logChan *chan log.Log, files map[string][2]string) {
	log.Delimiter(logChan)
	log.Msg(logChan, "changed files:")

	for path, hashes := range files {
		log.Msg(logChan, fmt.Sprintf("File %s has hashes:", path))

		for _, hash := range hashes {
			log.Msg(logChan, hash)
		}
	}
}

func printDups(logChan *chan log.Log, files map[string][]structs.FileInfo) {
	log.Delimiter(logChan)
	log.Msg(logChan, "dups:")

	for _, sameHashFiles := range files {
		if len(sameHashFiles) > 1 {
			log.Msg(logChan, fmt.Sprintf("Found dups for %v:", sameHashFiles[0].Path))

			for _, fileInfo := range sameHashFiles {
				log.Msg(logChan, fileInfo.Path)
			}
		}
	}
}
