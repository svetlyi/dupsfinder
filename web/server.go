package web

import (
	"github.com/svetlyi/dupsfinder/dups"
	"log"
	"net/http"
	"strconv"
)

func Serve(port int, stats *dups.Stats) {
	registerRoutes(stats)

	address := "localhost:" + strconv.Itoa(port)
	log.Printf("Web server address: %s\n", address)

	log.Fatal(http.ListenAndServe(address, nil))
}

func registerRoutes(stats *dups.Stats) {
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./web/templates/static"))))
	http.HandleFunc("/", mainpage(stats))
}
