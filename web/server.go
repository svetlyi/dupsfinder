package web

import (
	"github.com/svetlyi/dupsfinder/database"
	"github.com/svetlyi/dupsfinder/structs"
	"github.com/svetlyi/dupsfinder/web/templates/mainpage"
	"github.com/svetlyi/dupsfinder/web/templates/searchdups"
	"log"
	"net/http"
	"strconv"
)

func Serve(port int, stats *structs.Stats) {
	registerRoutes(stats)

	address := "localhost:" + strconv.Itoa(port)
	log.Printf("Web server address: %s\n", address)

	log.Fatal(http.ListenAndServe(address, nil))
}

func registerRoutes(stats *structs.Stats) {
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./web/templates/static"))))
	http.HandleFunc("/", mainpage.Mainpage(stats))
	http.HandleFunc("/search-dups", searchdups.Searchdups(database.GetDB()))
}
