package structs

import (
	"database/sql"
	"github.com/svetlyi/dupsfinder/log"
)

type App struct {
	DB      *sql.DB
	LogChan *chan log.Log
	// doneChan indicates that the application is about to exit
	// either user pressed exit or there was some fatal error
	ExitChan *chan bool
	// doneChan indicates that the calculations on all the files are finished
	DoneChan *chan bool
	Stats    *Stats
}

func NewApp() *App {
	logChan := make(chan log.Log)
	exitChan := make(chan bool)
	doneChan := make(chan bool)
	return &App{
		DB:       nil,
		LogChan:  &logChan,
		DoneChan: &doneChan,
		ExitChan: &exitChan,
		Stats:    &Stats{},
	}
}
