package searchdups

import (
	"database/sql"
	"github.com/svetlyi/dupsfinder/config"
	"github.com/svetlyi/dupsfinder/dups"
	"github.com/svetlyi/dupsfinder/structs"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type fileTmplObj struct {
	Path      string
	PathParts []string
	Hash      string
}

type searchDupsTmplObj struct {
	Files map[string][]fileTmplObj
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

	selectDupsStmt := dups.GetSelectDupsByDir(db)

	return func(writer http.ResponseWriter, request *http.Request) {
		dirsToSearch, dirsToSearchOk := request.URL.Query()["dir"]
		var page int = 1
		pages, pageOk := request.URL.Query()["page"]
		if pageOk {
			pageFromRequest, err := strconv.Atoi(pages[0])
			if nil == err {
				page = pageFromRequest
			}
		}
		var searchDupsObj searchDupsTmplObj
		searchDupsObj.Files = make(map[string][]fileTmplObj)

		if dirsToSearchOk && len(dirsToSearch) == 1 {
			rows, err := selectDupsStmt.Query(dirsToSearch[0], config.DupsPerPage, config.DupsPerPage*page)
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

				searchDupsObj.Files[fileInfo.Hash] = append(
					searchDupsObj.Files[fileInfo.Hash],
					fileTmplObj{
						Path:      fileInfo.Path,
						PathParts: fileInfo.SplitPath(),
						Hash:      fileInfo.Hash,
					},
				)
			}
			err = rows.Err()
			if err != nil {
				log.Fatal(err)
			}
		}
		rw := sync.RWMutex{}

		rw.RLock()
		mainTmplErr := mainTmpl.Execute(writer, searchDupsObj)
		if nil != mainTmplErr && !strings.Contains(mainTmplErr.Error(), "broken pipe") {
			log.Fatalf("Failed to execute a template. Error: %v", mainTmplErr)
		}
		rw.RUnlock()
	}
}
