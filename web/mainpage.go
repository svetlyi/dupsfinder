package web

import (
	"github.com/svetlyi/dupsfinder/dups"
	"html/template"
	"log"
	"net/http"
	"sync"
)

type mainpageObj struct {
	Stats string
}

func mainpage(stats *dups.Stats) func(writer http.ResponseWriter, request *http.Request) {
	mainTmpl := template.New("main.html")
	mainTmplFile := "web/templates/main.html"
	t, err := mainTmpl.ParseFiles(mainTmplFile)
	if nil != err {
		log.Fatalf("Failed to parse a template %s", mainTmplFile)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		rw := sync.RWMutex{}

		rw.RLock()
		err := t.Execute(writer, mainpageObj{stats.String()})
		if nil != err {
			log.Fatalf("Failed to execute a template %s. Error: %#v", mainTmplFile, err)
		}
		rw.RUnlock()
	}
}
