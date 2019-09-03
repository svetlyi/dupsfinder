package file

import (
	"crypto/md5"
	"fmt"
	"github.com/svetlyi/dupsfinder/dups"
	"io"
	"log"
	"os"
	"sync"
)

func ListenFilesChannel(filesChannel *chan string, filesInfoChannel *chan dups.FileInfo, wg *sync.WaitGroup) {
	for path := range *filesChannel {
		fileInfo, err := calculateHash(path)
		if nil == err {
			*filesInfoChannel <- fileInfo
		}
	}
	wg.Done()
}

func calculateHash(path string) (dups.FileInfo, error) {
	//defer wg.Done()
	f, err := os.Open(path)
	defer f.Close()

	if err != nil {
		log.Printf("Could not open the file: %q\n", path)
		return dups.FileInfo{}, err
	}

	h := md5.New()

	log.Printf("calculating hash for: %q\n", path)

	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	fileInfo := dups.FileInfo{
		Path: path,
		Hash: fmt.Sprintf("%x", h.Sum(nil)),
	}

	return fileInfo, nil
}
