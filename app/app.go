package app

import (
	"database/sql"
	"github.com/svetlyi/dupsfinder/config"
	"github.com/svetlyi/dupsfinder/logger"
	"github.com/svetlyi/dupsfinder/structs"
	"os"
)

type App struct {
	DB     *sql.DB
	Logger *logger.Logger
	// doneChan indicates that the application is about to exit
	// either user pressed exit or there was some fatal error
	ExitChan *structs.ExitChan
	// doneChan indicates that the calculations on all the files are finished
	DoneChan *chan bool
	Stats    *structs.Stats
}

// Fatal logs some really fatal error and tries to finish
// the application gracefully by closing ExitChan channel
func (app *App) Fatal(msg string) {
	app.Logger.Err(msg)
	close(*app.ExitChan)
}

func New() (*App, error) {
	exitChan := make(structs.ExitChan)
	doneChan := make(chan bool)
	if logWriter, err := os.Create(config.LogFile); nil == err {
		return &App{
			DB:       nil,
			Logger:   logger.New(logWriter),
			DoneChan: &doneChan,
			ExitChan: &exitChan,
			Stats:    &structs.Stats{},
		}, nil
	} else {
		return &App{}, err
	}
}
