package file

import (
	"fmt"
	"github.com/svetlyi/dupsfinder/app"
	"github.com/svetlyi/dupsfinder/config"
	"github.com/svetlyi/dupsfinder/structs"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// walkThroughFiles walks through the files in initDir folder
// and sends the files to filesChannel channel
func walkThroughFiles(initDir string, filesChan *chan structs.FileInfo, app *app.App) {
	mutex := sync.Mutex{}
	defer close(*filesChan)

	err := filepath.Walk(initDir, func(path string, info os.FileInfo, err error) error {
		select {
		case <-*app.ExitChan:
			return io.EOF
		default:
			if err != nil {
				app.Logger.Err(fmt.Sprintf("file.WalkThroughFiles: error accessing a path %q: %v\n", path, err))
				return err
			}
			app.Logger.Msg(fmt.Sprintf("visited file or dir: %q\n", path))
			if !info.IsDir() && (info.Size() > config.IgnoreFilesLessThanBytes) {
				app.Logger.Msg(fmt.Sprintf("it is not a dir: %q\n", path))

				mutex.Lock()
				(*app.Stats).FilesAmount++
				(*app.Stats).FilesSize += info.Size()
				mutex.Unlock()

				*filesChan <- structs.FileInfo{Path: path, Hash: "", Size: info.Size()}
			}
			return nil
		}
	})

	if nil != err && io.EOF != err {
		app.Logger.Msg(fmt.Sprintf("error walking the path: %s\n", initDir))
	} else if io.EOF == err {
		app.Logger.Msg("file walking was terminated by a user")
	}
}
