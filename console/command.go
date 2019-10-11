package console

import (
	"bufio"
	"fmt"
	"github.com/svetlyi/dupsfinder/app"
	"github.com/svetlyi/dupsfinder/config"
	"github.com/svetlyi/dupsfinder/dups"
	"github.com/svetlyi/dupsfinder/file"
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
		fmt.Sprintf("choose command (%s)", strings.Join(getCommandList(), ", ")),
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
func runExitCmd(app *app.App) {
	select {
	case <-*app.ExitChan:
	default:
		close(*app.ExitChan)
		fmt.Println("closed exit channel. It will exit as soon as the operations will be finished")

		go func(exitChan *chan bool) {
			<-*exitChan
			fmt.Println("exiting...")
			os.Exit(0)
		}(app.ExitChan)
	}
}

func runWebServer(scanner *bufio.Scanner, app *app.App) {
	portText := getValueFromUser(scanner, "select port for the server", strconv.Itoa(config.WebServerPort))

	port64, err := strconv.ParseUint(portText, 10, 16)
	if nil != err {
		fmt.Printf("wrong value for port: %s\n", err.Error())
	} else {
		go web.Serve(uint16(port64), app)
		fmt.Println("web server is about to run")
	}
}

func runUpdateIndexDBCmd(scanner *bufio.Scanner, app *app.App) {
	app.Stats.StartTime = time.Now()

	var path string

	path = getValueFromUser(scanner, "type folder to update index from", "")

	procNumText := getValueFromUser(scanner, "type number of goroutines (processors) to use", strconv.Itoa(config.ProcNum))
	procNum64, err := strconv.ParseInt(procNumText, 10, 8)
	if nil != err {
		fmt.Printf("either wrong number of goroutines or wrong format: %s\n", err)
	}
	if path != "" && procNum64 > 0 {
		go file.UpdateIndexDB(path, uint8(procNum64), app)
	}
}

func runShowStatsCmd(stats *structs.Stats) {
	fmt.Println(stats.String())
}

func runShowDupsCmd(scanner *bufio.Scanner, app *app.App) {
	var path string
	var page int

	path = getValueFromUser(scanner, "type folder where to show dups from", "")

	var enoughPages = false

	for !enoughPages {
		duplicatesSets, err := dups.Get(path, page, app)
		if err != nil {
			app.Logger.Err(err.Error())
			return
		}

		if len(duplicatesSets.Files) > 0 {
			app.Logger.Msg(fmt.Sprintf("found %d duplicates in %s", len(duplicatesSets.Files), path))

			for hash, duplicates := range duplicatesSets.Files {
				fmt.Printf("\n=======\n%s:\n=======\n\n", hash)
				for _, duplicate := range duplicates {
					fmt.Println(duplicate.Path)
				}
			}
			page++
		} else {
			fmt.Println("there are no more duplicates")
			break
		}
		moreValues := getValueFromUser(scanner, "more values? (1/0)", "1")
		enoughPages = enoughPages || ("0" == moreValues)
	}
}

func runWrongCmd(cmd string) {
	similar := findSimilar(cmd, getCommandList())
	fmt.Printf("wrong command: %s. maybe you meant '%s'?\n", cmd, similar)
}
