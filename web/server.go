package web

import (
	"fmt"
	"github.com/svetlyi/dupsfinder/log"
	"github.com/svetlyi/dupsfinder/structs"
	"github.com/svetlyi/dupsfinder/web/templates/mainpage"
	"github.com/svetlyi/dupsfinder/web/templates/searchdups"
	"net/http"
)

func Serve(port uint16, app *structs.App) {
	registerRoutes(app)

	address := fmt.Sprintf("localhost:%d", port)
	log.Msg(app.LogChan, fmt.Sprintf("web server address: %s\n", address))

	if err := http.ListenAndServe(address, nil); nil != err {
		log.Err(app.LogChan, err.Error())
	}
}

func registerRoutes(app *structs.App) {
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./web/templates/static"))))
	http.HandleFunc("/", mainpage.Mainpage(app.Stats))
	http.HandleFunc("/search-dups", searchdups.Searchdups(app.DB))
}
