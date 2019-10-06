package file

import (
	"crypto/md5"
	"fmt"
	"github.com/svetlyi/dupsfinder/logger"
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
func CalculateHashes(filesChan *chan string, filesInfoChan *chan structs.FileInfo, calcHashesWG *sync.WaitGroup, logger *logger.Logger) {
	var hash string

	for path := range *filesChan {
		hash = calculateHash(path, logger)
		fileInfo := structs.FileInfo{
			Path: path,
			Hash: hash,
		}
		*filesInfoChan <- fileInfo
	}
	calcHashesWG.Done()
}

func calculateHash(path string, logger *logger.Logger) string {
	f, err := os.Open(path)

	if err != nil {
		logger.Err(fmt.Sprintf("could not open the file: %q", path))
	}
	defer f.Close()

	h := md5.New()

	logger.Msg(fmt.Sprintf("calculating hash for: %q", path))

	if _, err := io.Copy(h, f); err != nil {
		logger.Err(err.Error())
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
