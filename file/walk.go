package file

import (
	"fmt"
	"github.com/svetlyi/dupsfinder/log"
	"github.com/svetlyi/dupsfinder/structs"
	"io"
	"os"
	"path/filepath"
	"sync"
)

/**
Walks through the files in initDir folder
and sends the files to filesChannel channel
*/
func WalkThroughFiles(initDir string, filesChannel *chan string, app *structs.App) {
	mutex := sync.Mutex{}
	defer close(*filesChannel)

	err := filepath.Walk(initDir, func(path string, info os.FileInfo, err error) error {
		select {
		case <-*app.ExitChan:
			return io.EOF
		default:
			if err != nil {
				log.Err(app.LogChan, fmt.Sprintf("file.WalkThroughFiles: error accessing a path %q: %v\n", path, err))
				return err
			}
			log.Msg(app.LogChan, fmt.Sprintf("visited file or dir: %q\n", path))
			if !info.IsDir() {
				log.Msg(app.LogChan, fmt.Sprintf("it is not a dir: %q\n", path))

				mutex.Lock()
				(*app.Stats).FilesAmount++
				(*app.Stats).FilesSize += info.Size()
				mutex.Unlock()
				*filesChannel <- path
			}
			return nil
		}
	})

	if nil != err && io.EOF != err {
		log.Err(app.LogChan, fmt.Sprintf("error walking the path: %s\n", initDir))
	} else if io.EOF == err {
		log.Msg(app.LogChan, "file walking was terminated by a user")
	}
}
