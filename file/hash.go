package file

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"github.com/svetlyi/dupsfinder/config"
	"github.com/svetlyi/dupsfinder/structs"
	"io"
	"log"
	"os"
	"sync"
)

/**
Listens to filesChannel channel and calculates hashes for
the received from the channel files. Information about the
processed files goes into filesInfoChannel channel.
*/
func ListenFilesChannel(filesChannel *chan string, filesInfoChannel *chan structs.FileInfo, wg *sync.WaitGroup, db *sql.DB) {
	selectStmt := GetSelectHashByPathStmt(db)
	defer selectStmt.Close()
	var hash string
	var selectHashErr error

	for path := range *filesChannel {
		if config.CheckFilesIntegrity {
			hash = calculateHash(path)
		} else {
			hash, selectHashErr = GetHashByPathFromDB(selectStmt, path)
			if selectHashErr == nil {
				log.Printf("Found file %s in database\n", path)
			} else if selectHashErr == sql.ErrNoRows {
				hash = calculateHash(path)
			}
		}
		fileInfo := structs.FileInfo{
			Path: path,
			Hash: hash,
		}
		*filesInfoChannel <- fileInfo
	}
	wg.Done()
}

func calculateHash(path string) string {
	f, err := os.Open(path)
	defer f.Close()

	if err != nil {
		log.Fatalf("Could not open the file: %q\n", path)
	}

	h := md5.New()

	log.Printf("calculating hash for: %q\n", path)

	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
