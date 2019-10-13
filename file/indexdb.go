package file

import (
	"database/sql"
	"fmt"
	"github.com/svetlyi/dupsfinder/app"
	"github.com/svetlyi/dupsfinder/structs"
	"sync"
)

// updates database index, populates it with new files and
// removes old ones
func UpdateIndexDB(path string, procNum uint8, app *app.App) {
	var filesChan = make(chan structs.FileInfo)
	var filesHashesChan = make(chan structs.FileInfo)

	go WalkThroughFiles(path, &filesChan, app)

	var wgIndex uint8
	calculateHashesWG := sync.WaitGroup{}
	for wgIndex = 1; wgIndex <= procNum; wgIndex++ {
		app.Logger.Msg(fmt.Sprintf("creating a go routine %d to calculate hashes\n", wgIndex))
		calculateHashesWG.Add(1)
		go CalculateHashes(&filesChan, &filesHashesChan, &calculateHashesWG, app.Logger)
	}

	go storeFilesInfo(&filesHashesChan, app)

	calculateHashesWG.Wait()
	app.Logger.Msg("closing files info channel")
	close(filesHashesChan)
}

// storeFilesInfo listens to filesInfoChan and stores
// the information in a database
func storeFilesInfo(filesInfoChan *chan structs.FileInfo, app *app.App) {
	selectStmt := GetSelectHashByPathStmt(app.DB)
	defer selectStmt.Close()
	insertStmt := GetInsertStmt(app.DB)
	defer insertStmt.Close()
	updateHashStmt := GetUpdateHashStmt(app.DB)
	defer updateHashStmt.Close()

	for fileInfo := range *filesInfoChan {
		hashInDB, hashInDBErr := GetHashByPathFromDB(selectStmt, fileInfo.Path)
		if nil == hashInDBErr {
			app.Logger.Msg(fmt.Sprintf("File %s already exists in DB\n", fileInfo.Path))
			if hashInDB != fileInfo.Hash {
				app.Logger.Msg(fmt.Sprintf("File %s has been changed\n", fileInfo.Path))
				if _, err := updateHashStmt.Exec(fileInfo.Hash, fileInfo.Path); nil != err {
					app.Fatal(err.Error())
				}
			}
		} else if hashInDBErr == sql.ErrNoRows {
			if _, err := insertStmt.Exec(fileInfo.Path, fileInfo.Hash, fileInfo.Size); nil != err {
				app.Fatal(err.Error())
			}
		} else {
			app.Fatal(hashInDBErr.Error())
		}
	}

	*app.DoneChan <- true
}
