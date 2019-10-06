package web

import (
	"context"
	"fmt"
	"github.com/svetlyi/dupsfinder/app"
	"github.com/svetlyi/dupsfinder/web/templates/mainpage"
	"github.com/svetlyi/dupsfinder/web/templates/searchdups"
	"net/http"
)

func Serve(port uint16, app *app.App) {
	registerRoutes(app)

	address := fmt.Sprintf("localhost:%d", port)
	srv := http.Server{Addr: address}
	app.Logger.Msg(fmt.Sprintf("web server address: %s\n", address))

	go func(server http.Server) {
		if err := srv.ListenAndServe(); nil != err {
			app.Logger.Err(err.Error())
		}
	}(srv)
	<-*app.ExitChan
	if err := srv.Shutdown(context.Background()); nil != err {
		fmt.Printf("server couldn't terminate gracefully: %s\n", err.Error())
	} else {
		fmt.Println("\nserver terminated successfully")
	}
}

func registerRoutes(app *app.App) {
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./web/templates/static"))))
	http.HandleFunc("/", mainpage.Mainpage(app.Stats))
	http.HandleFunc("/search-dups", searchdups.Searchdups(app.DB))
}
