package console

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/svetlyi/dupsfinder/app"
	"github.com/svetlyi/dupsfinder/config"
	"github.com/svetlyi/dupsfinder/database"
	"os"
)

var application *app.App

func init() {
	var err error
	if application, err = app.New(); nil != err {
		fmt.Println(err)
		os.Exit(1)
	} else {
		go application.Logger.ListenToChannel(application.ExitChan)
	}
}

func Run() {
	scanner := bufio.NewScanner(os.Stdin)
	var cmd string

	for {
		cmd = ""
		if nil == application.DB {
			application.DB = runCreateDB(scanner)
		} else if "" == cmd {
			cmd = chooseCmd(scanner)
			runCmd(cmd, scanner)
		}
	}
}

func runCreateDB(scanner *bufio.Scanner) *sql.DB {
	dbPath := getValueFromUser(scanner, fmt.Sprintf("choose db path"), config.DBPath)

	return database.NewDB(dbPath)
}

func runCmd(cmd string, scanner *bufio.Scanner) {
	switch cmd {
	case findDupsCmd:
		runFindDupsCmd(scanner, application)
	case showStatsCmd:
		runShowStatsCmd(application.Stats)
	case runWebServerCmd:
		runWebServer(scanner, application)
	case exitCmd:
		runExitCmd(application)
	default:
		runWrongCmd(cmd)
		return
	}
}
