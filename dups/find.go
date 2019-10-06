package dups

import (
	"fmt"
	"github.com/svetlyi/dupsfinder/app"
	"github.com/svetlyi/dupsfinder/file"
	"github.com/svetlyi/dupsfinder/structs"
	"sync"
)

func Find(path string, procNum uint8, app *app.App) {
	var filesChan = make(chan string)
	var filesInfoChan = make(chan structs.FileInfo)

	go file.WalkThroughFiles(path, &filesChan, app)

	var wgIndex uint8
	calculateHashesWG := sync.WaitGroup{}
	for wgIndex = 1; wgIndex <= procNum; wgIndex++ {
		app.Logger.Msg(fmt.Sprintf("creating a go routine %d to calculate hashes\n", wgIndex))
		calculateHashesWG.Add(1)
		go file.CalculateHashes(&filesChan, &filesInfoChan, &calculateHashesWG, app.Logger)
	}

	go storeFilesInfo(&filesInfoChan, app)

	calculateHashesWG.Wait()
	app.Logger.Msg("closing files info channel")
	close(filesInfoChan)
}
