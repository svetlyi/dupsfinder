package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/svetlyi/dupsfinder/database"
	"github.com/svetlyi/dupsfinder/dups"
	"github.com/svetlyi/dupsfinder/file"
	"github.com/svetlyi/dupsfinder/web"
	"log"
	"os"
	"sync"
	"time"
)

var stats *dups.Stats = &dups.Stats{}

var path string
var procNum int
var dbPath string
var port int

func init() {
	flag.StringVar(&path, "path", ".", "a path where to look for duplicates")
	flag.StringVar(&dbPath, "dbPath", "", "a path where the database would be stored")
	flag.IntVar(&procNum, "procNum", 1, "a number of processors to use")
	flag.IntVar(&port, "port", 55786, "a web server port on which statistics can be shown")
}

func main() {
	checkArgs()
	database.CreateDB(dbPath)

	var filesChannel chan string = make(chan string)
	var filesInfoChannel chan dups.FileInfo = make(chan dups.FileInfo)
	// doneChannel indicates that the calculations on all the files are finished
	var doneChannel chan bool = make(chan bool)

	stats.StartTime = time.Now()
	go web.Serve(port, stats)

	go file.WalkThroughFiles(path, &filesChannel, stats)

	wg := sync.WaitGroup{}
	var wgIndex int
	for wgIndex = 1; wgIndex <= procNum; wgIndex++ {
		log.Printf("creating a go routine %d to calculate hashes\n", wgIndex)
		wg.Add(1)
		go file.ListenFilesChannel(&filesChannel, &filesInfoChannel, &wg, database.GetDB())
	}

	go dups.ListenFilesInfoChannel(&filesInfoChannel, &doneChannel, database.GetDB())

	wg.Wait()
	fmt.Println("closing files info channel")
	close(filesInfoChannel)

	stats.EndTime = time.Now()
	log.Println(stats.String())

	<-doneChannel
	confirmExit()
}

func confirmExit() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Press enter to exit: ")
	_, err := reader.ReadString('\n')
	if nil != err {
		log.Fatal("An error reading stdin: ", err)
	}
}

func checkArgs() {
	flag.Parse()

	if procNum < 1 {
		log.Fatalf("Wrong number of processors specified: %d", procNum)
		return
	}
	if port > 65535 || port < 1 {
		log.Fatalf("Wrong port number: %d", port)
		return
	}
	if "" == dbPath {
		log.Fatal("Specify path where to store the database!")
		return
	}
}
