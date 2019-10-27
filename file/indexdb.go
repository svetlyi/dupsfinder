package file

import (
	"database/sql"
	"fmt"
	"github.com/svetlyi/dupsfinder/app"
	"github.com/svetlyi/dupsfinder/structs"
	"os"
	"sync"
	"time"
)

// UpdateIndexDB updates database index, populates it
// with new files and removes old ones
func UpdateIndexDB(path string, procNum uint8, app *app.App) {
	var filesChan = make(chan structs.FileInfo)
	var filesHashesChan = make(chan structs.FileInfo)

	go walkThroughFiles(path, &filesChan, app)

	var wgIndex uint8
	calculateHashesWG := sync.WaitGroup{}
	for wgIndex = 1; wgIndex <= procNum; wgIndex++ {
		app.Logger.Msg(fmt.Sprintf("creating a go routine %d to calculate hashes\n", wgIndex))
		calculateHashesWG.Add(1)
		go calculateHashes(&filesChan, &filesHashesChan, &calculateHashesWG, app.Logger)
	}

	go storeFilesInfo(&filesHashesChan, app)

	calculateHashesWG.Wait()
	app.Logger.Msg("closing files info channel")
	close(filesHashesChan)
	removeNonExistingFilesFromDb(app)
	var mutex sync.Mutex = sync.Mutex{}
	mutex.Lock()
	(*app.Stats).EndTime = time.Now()
	mutex.Unlock()
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

func removeNonExistingFilesFromDb(app *app.App) {
	app.Logger.Msg("removing files...")
	stmt, err := app.DB.Prepare("SELECT path FROM files LIMIT ? OFFSET ?")
	if err != nil {
		app.Fatal("removeNonExistingFilesFromDb could not prepare statement: " + err.Error())
	}
	defer stmt.Close()

	var page int
	var morePages = true

	for morePages {
		morePages = false
		removeBatchNonExistingFiles(page, &morePages, stmt, app)
		page++

		select {
		case <-*app.ExitChan:
			return
		default:
			continue
		}
	}
}

func removeBatchNonExistingFiles(page int, morePages *bool, stmt *sql.Stmt, app *app.App) {
	const rowsLimit = 100

	rows, err := stmt.Query(rowsLimit, rowsLimit*page)
	if err != nil {
		app.Fatal("removeNonExistingFilesFromDb could not fetch rows: " + err.Error())
	}
	var path string
	var pathsToDelete [rowsLimit]string
	var pathsToDeleteCounter int

	defer rows.Close()
	for rows.Next() {
		*morePages = true
		err := rows.Scan(&path)
		if err != nil {
			app.Fatal("removeNonExistingFilesFromDb could not fetch rows: " + err.Error())
			break
		}
		if !fileExists(path) {
			app.Logger.Msg(fmt.Sprintf("adding file %s to delete", path))
			pathsToDelete[pathsToDeleteCounter] = path
			pathsToDeleteCounter++
		}
	}
	if err = rows.Err(); err != nil {
		app.Fatal("An unexpected error: " + err.Error())
	}

	for i := 0; i < pathsToDeleteCounter; i++ {
		if "" == pathsToDelete[i] {
			continue
		}
		app.Logger.Msg(fmt.Sprintf("removing file '%s' from database", pathsToDelete[i]))
		_, err = app.DB.Exec("DELETE FROM files WHERE path = ?", pathsToDelete[i])
		if err != nil {
			app.Fatal(fmt.Sprintf("failed to remove file %s. Error: %s", pathsToDelete[i], err.Error()))
			break
		}
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
