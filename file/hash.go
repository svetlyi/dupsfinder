package file

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"github.com/svetlyi/dupsfinder/dups"
	"io"
	"log"
	"os"
	"sync"
)

func ListenFilesChannel(filesChannel *chan string, filesInfoChannel *chan dups.FileInfo, wg *sync.WaitGroup, db *sql.DB) {
	selectStmt := getSelectHashByPathStmt(db)
	defer selectStmt.Close()
	var hash string

	for path := range *filesChannel {
		selectErr := selectStmt.QueryRow(path).Scan(&hash)
		if selectErr != nil {
			if selectErr == sql.ErrNoRows {
				var hashErr error
				hash, hashErr = calculateHash(path)
				if nil != hashErr {
					log.Fatal(hashErr)
				}
			} else {
				log.Fatal(selectErr)
			}
		} else {
			log.Printf("Found file %s in database\n", path)
		}
		fileInfo := dups.FileInfo{
			Path: path,
			Hash: hash,
		}
		*filesInfoChannel <- fileInfo
	}
	wg.Done()
}

func calculateHash(path string) (string, error) {
	//defer wg.Done()
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

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
