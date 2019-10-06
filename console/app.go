package console

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/svetlyi/dupsfinder/config"
	"github.com/svetlyi/dupsfinder/database"
	"github.com/svetlyi/dupsfinder/log"
	"github.com/svetlyi/dupsfinder/structs"
	"os"
)

var app *structs.App = structs.NewApp()

func init() {
	go log.ListenToChannel(app.LogChan)
}

func Run() {
	scanner := bufio.NewScanner(os.Stdin)
	var cmd string

	for {
		cmd = ""
		if nil == app.DB {
			app.DB = runCreateDB(scanner)
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
		runFindDupsCmd(scanner, app)
	case showStatsCmd:
		runShowStatsCmd(app.Stats)
	case showLastLogsCmd:
		runShowLastLogs(scanner)
	case runWebServerCmd:
		runWebServer(scanner, app)
	case exitCmd:
		runExitCmd(app)
	default:
		runWrongCmd(cmd)
		return
	}
}
