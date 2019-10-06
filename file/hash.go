package file

import (
	"crypto/md5"
	"fmt"
	"github.com/svetlyi/dupsfinder/log"
	"github.com/svetlyi/dupsfinder/structs"
	"io"
	"os"
	"sync"
)

/**
Listens to filesChannel channel and calculates hashes for
the received from the channel files. Information about the
processed files goes into filesInfoChannel channel.
*/
func CalculateHashes(filesChan *chan string, filesInfoChan *chan structs.FileInfo, calcHashesWG *sync.WaitGroup, logChan *chan log.Log) {
	var hash string

	for path := range *filesChan {
		hash = calculateHash(path, logChan)
		fileInfo := structs.FileInfo{
			Path: path,
			Hash: hash,
		}
		*filesInfoChan <- fileInfo
	}
	calcHashesWG.Done()
}

func calculateHash(path string, logChan *chan log.Log) string {
	f, err := os.Open(path)

	if err != nil {
		log.Err(logChan, fmt.Sprintf("could not open the file: %q\n", path))
	}
	defer f.Close()

	h := md5.New()

	log.Msg(logChan, fmt.Sprintf("calculating hash for: %q\n", path))

	if _, err := io.Copy(h, f); err != nil {
		log.Err(logChan, err.Error())
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
