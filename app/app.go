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
	ExitChan *chan bool
	// doneChan indicates that the calculations on all the files are finished
	DoneChan *chan bool
	Stats    *structs.Stats
}

func New() (*App, error) {
	exitChan := make(chan bool)
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
