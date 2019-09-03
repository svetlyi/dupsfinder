package file

import (
	"fmt"
	"github.com/svetlyi/dupsfinder/dups"
	"log"
	"os"
	"path/filepath"
)

var stats dups.Stats = dups.Stats{}

func WalkThroughFiles(initDir string, filesChannel *chan string) error {
	defer close(*filesChannel)

	err := filepath.Walk(initDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("error accessing a path %q: %v\n", path, err)
			return err
		}
		log.Printf("visited file or dir: %q\n", path)
		if !info.IsDir() {
			log.Printf("it is not a dir: %q\n", path)

			stats.FilesAmount++
			stats.FilesSize += info.Size()
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
