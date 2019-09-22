package mainpage

import (
	"github.com/svetlyi/dupsfinder/structs"
	"html/template"
	"log"
	"net/http"
	"sync"
)

type mainpageObj struct {
	Stats string
}

func Mainpage(stats *structs.Stats) func(writer http.ResponseWriter, request *http.Request) {
	var templatesPaths []string = []string{
		"web/templates/layout/header.html",
		"web/templates/layout/footer.html",
		"web/templates/mainpage/content.html",
	}
	templates, err := template.ParseFiles(templatesPaths...)
	if nil != err {
		log.Fatalf("Failed to parse templatesPaths %v: %v", templatesPaths, err)
	}
	mainTmpl := templates.Lookup("content")

	return func(writer http.ResponseWriter, request *http.Request) {
		rw := sync.RWMutex{}

		rw.RLock()
		err := mainTmpl.Execute(writer, mainpageObj{stats.String()})
		if nil != err {
			log.Fatalf("Failed to execute a template %#v. Error: %#v", mainTmpl, err)
		}
		rw.RUnlock()
	}
}
