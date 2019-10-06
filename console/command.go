package console

import (
	"bufio"
	"fmt"
	"github.com/svetlyi/dupsfinder/config"
	"github.com/svetlyi/dupsfinder/dups"
	"github.com/svetlyi/dupsfinder/log"
	"github.com/svetlyi/dupsfinder/structs"
	"github.com/svetlyi/dupsfinder/web"
	"os"
	"strconv"
	"strings"
	"time"
)

func chooseCmd(scanner *bufio.Scanner) string {
	cmd := getValueFromUser(
		scanner,
		fmt.Sprintf("choose command (%s): ", strings.Join(getCommandList(), ", ")),
		"",
	)

	return cmd
}

func getValueFromUser(scanner *bufio.Scanner, text string, def string) string {
	if "" != def {
		fmt.Printf("%s %s (%s): ", basePrompt, text, def)
	} else {
		fmt.Printf("%s %s: ", basePrompt, text)
	}
	scanner.Scan()
	text = scanner.Text()
	if "" == text && "" != def {
		text = def
	}
	return text
}

// gracefully terminates the application
// listens to the done channel and if it is closed,
// it exits
func runExitCmd(app *structs.App) {
	select {
	case <-*app.ExitChan:
	default:
		close(*app.ExitChan)
		fmt.Println("closed exit channel. It will exit as soon as the operations will be finished")

		go func(doneChan *chan bool) {
			<-*doneChan
			fmt.Println("exiting...")
			os.Exit(0)
		}(app.DoneChan)
	}
}

func runWebServer(scanner *bufio.Scanner, app *structs.App) {
	portText := getValueFromUser(scanner, "select port for the server", strconv.Itoa(config.WebServerPort))

	port64, err := strconv.ParseUint(portText, 10, 16)
	if nil != err {
		fmt.Printf("wrong value for port: %s\n", err.Error())
	} else {
		go web.Serve(uint16(port64), app)
	}
}

func runFindDupsCmd(scanner *bufio.Scanner, app *structs.App) {
	app.Stats.StartTime = time.Now()

	var path string

	path = getValueFromUser(scanner, "type folder to search dups in", "")

	procNumText := getValueFromUser(scanner, "type number of goroutines (processors) to use", strconv.Itoa(config.ProcNum))
	procNum64, err := strconv.ParseInt(procNumText, 10, 8)
	if nil != err {
		fmt.Printf("either wrong number of goroutines or wrong format: %s\n", err)
	}
	if path != "" && procNum64 > 0 {
		go dups.Find(path, uint8(procNum64), app)
	}
}

func runShowLastLogs(scanner *bufio.Scanner) {
	numOfMessages := getValueFromUser(scanner, "number of messages", "10")
	num, err := strconv.Atoi(numOfMessages)
	if nil != err {
		fmt.Printf("wrong number of messages: %s", err.Error())
		num = 10
	}
	for _, msg := range log.GetLastMessages(num) {
		fmt.Println(msg)
	}
}

func runShowStatsCmd(stats *structs.Stats) {
	fmt.Println(stats.String())
}

func runWrongCmd(cmd string) {
	similar := findSimilar(cmd, getCommandList())
	fmt.Printf("wrong command: %s. maybe you meant '%s'?\n", cmd, similar)
}
