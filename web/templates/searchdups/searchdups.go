package searchdups

import (
	"database/sql"
	"github.com/svetlyi/dupsfinder/file"
	"github.com/svetlyi/dupsfinder/structs"
	"html/template"
	"log"
	"net/http"
	"sync"
)

type searchDupsObj struct {
	Files []structs.FileInfo
}

func Searchdups(db *sql.DB) func(writer http.ResponseWriter, request *http.Request) {
	var templatesPaths []string = []string{
		"web/templates/layout/header.html",
		"web/templates/layout/error.html",
		"web/templates/layout/footer.html",
		"web/templates/searchdups/content.html",
	}
	templates, err := template.ParseFiles(templatesPaths...)
	if nil != err {
		log.Fatalf("Failed to parse templatesPaths %v: %v", templatesPaths, err)
	}
	mainTmpl := templates.Lookup("content")

	selectFilesStmt := file.GetSelectFilesByDir(db)

	return func(writer http.ResponseWriter, request *http.Request) {
		dirsToSearch, ok := request.URL.Query()["dir"]
		var searchDups searchDupsObj

		if ok && len(dirsToSearch) == 1 {
			rows, err := selectFilesStmt.Query(dirsToSearch[0] + "%")
			if nil != err {
				log.Fatalf("An error while querying: %v", err)
			}
			defer rows.Close()
			for rows.Next() {
				var fileInfo structs.FileInfo
				err := rows.Scan(&fileInfo.Path, &fileInfo.Hash)
				if err != nil {
					log.Fatal(err)
				}
				searchDups.Files = append(searchDups.Files, fileInfo)
			}
			err = rows.Err()
			if err != nil {
				log.Fatal(err)
			}
		}
		rw := sync.RWMutex{}

		rw.RLock()
		mainTmplErr := mainTmpl.Execute(writer, searchDups)
		if nil != mainTmplErr {
			log.Fatalf("Failed to execute a template %s. Error: %#v", mainTmpl, mainTmplErr)
		}
		rw.RUnlock()
	}
}
