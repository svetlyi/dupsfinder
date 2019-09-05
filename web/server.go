package web

import (
	"fmt"
	"github.com/svetlyi/dupsfinder/dups"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func Serve(port int, stats *dups.Stats) {
	rw := sync.RWMutex{}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		rw.RLock()
		fmt.Fprintf(writer, stats.String())
		rw.RUnlock()
	})
	address := "localhost:" + strconv.Itoa(port)
	log.Printf("Web server address: %s\n", address)

	log.Fatal(http.ListenAndServe(address, nil))
}
