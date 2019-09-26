package file

import (
	"fmt"
	"github.com/svetlyi/dupsfinder/structs"
	"log"
	"os"
	"path/filepath"
	"sync"
)

/**
Walks through the files in initDir folder
and sends the files to filesChannel channel
*/
func WalkThroughFiles(initDir string, filesChannel *chan string, stats *structs.Stats) error {
	mutex := sync.Mutex{}
	defer close(*filesChannel)

	err := filepath.Walk(initDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("file.WalkThroughFiles: error accessing a path %q: %v\n", path, err)
			return err
		}
		log.Printf("visited file or dir: %q\n", path)
		if !info.IsDir() {
			log.Printf("it is not a dir: %q\n", path)

			mutex.Lock()
			(*stats).FilesAmount++
			(*stats).FilesSize += info.Size()
			mutex.Unlock()
			*filesChannel <- path
		}
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path: %v\n", initDir)
		return err
	}

	return nil
}
