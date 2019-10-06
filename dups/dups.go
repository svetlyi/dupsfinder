package dups

import (
	"database/sql"
	"fmt"
	"github.com/svetlyi/dupsfinder/app"
	"github.com/svetlyi/dupsfinder/file"
	"github.com/svetlyi/dupsfinder/logger"
	"github.com/svetlyi/dupsfinder/structs"
)

// storeFilesInfo listens to filesInfoChannel and stores
// the information in a databse
func storeFilesInfo(filesInfoChannel *chan structs.FileInfo, app *app.App) {
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
			app.Logger.Msg(fmt.Sprintf("File %s already exists in DB\n", fileInfo.Path))
			if hashInDB != fileInfo.Hash {
				app.Logger.Msg(fmt.Sprintf("File %s was changed\n", fileInfo.Path))
				filesWithChangedHashes[fileInfo.Path] = [2]string{fileInfo.Hash, hashInDB}
			}
		} else if hashInDBErr == sql.ErrNoRows {
			if _, err := insertStmt.Exec(fileInfo.Path, fileInfo.Hash); nil != err {
				app.Logger.Err(err.Error())
				close(*app.ExitChan)
			}
		} else {
			app.Logger.Err(hashInDBErr.Error())
			close(*app.ExitChan)
		}
		if nil == fileDups[fileInfo.Hash] {
			fileDups[fileInfo.Hash] = make([]structs.FileInfo, 0)
		}
		fileDups[fileInfo.Hash] = append(fileDups[fileInfo.Hash], fileInfo)
	}

	printDups(app.Logger, fileDups)
	printChangedFiles(app.Logger, filesWithChangedHashes)
	*app.DoneChan <- true
}

func printChangedFiles(logger *logger.Logger, files map[string][2]string) {
	logger.Delimiter()
	logger.Msg("changed files:")

	for path, hashes := range files {
		logger.Msg(fmt.Sprintf("File %s has hashes:", path))

		for _, hash := range hashes {
			logger.Msg(hash)
		}
	}
}

func printDups(logger *logger.Logger, files map[string][]structs.FileInfo) {
	logger.Delimiter()
	logger.Msg("dups:")

	for _, sameHashFiles := range files {
		if len(sameHashFiles) > 1 {
			logger.Msg(fmt.Sprintf("Found dups for %v:", sameHashFiles[0].Path))

			for _, fileInfo := range sameHashFiles {
				logger.Msg(fileInfo.Path)
			}
		}
	}
}
