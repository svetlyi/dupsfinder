package dups

import (
	"fmt"
	"github.com/svetlyi/dupsfinder/app"
	"github.com/svetlyi/dupsfinder/dups"
	"github.com/svetlyi/dupsfinder/structs"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func Show(app *app.App) func(writer http.ResponseWriter, request *http.Request) {
	var templatesPaths []string = []string{
		"web/templates/layout/header.html",
		"web/templates/layout/error.html",
		"web/templates/layout/footer.html",
		"web/templates/dups/content.html",
	}
	templates, err := template.ParseFiles(templatesPaths...)
	if nil != err {
		log.Fatalf("Failed to parse templatesPaths %v: %v", templatesPaths, err)
	}
	mainTmpl := templates.Lookup("content")

	return func(writer http.ResponseWriter, request *http.Request) {
		dirsToSearch, dirsToSearchOk := request.URL.Query()["dir"]
		var page int = 0
		pages, pageOk := request.URL.Query()["page"]
		if pageOk {
			pageFromRequest, err := strconv.Atoi(pages[0])
			if nil == err {
				page = pageFromRequest
			}
		}

		var searchDupsObj structs.DupsTmplObj

		if dirsToSearchOk && len(dirsToSearch) == 1 {
			searchDupsObj, err = dups.Get(dirsToSearch[0], page, app)
			if nil != err {
				searchDupsObj = structs.DupsTmplObj{}
			}
		}
		rw := sync.RWMutex{}

		rw.RLock()
		mainTmplErr := mainTmpl.Execute(writer, searchDupsObj)
		if nil != mainTmplErr && !strings.Contains(mainTmplErr.Error(), "broken pipe") {
			app.Logger.Err(fmt.Sprintf("Failed to execute a template. Error: %v", mainTmplErr))
		}
		rw.RUnlock()
	}
}
