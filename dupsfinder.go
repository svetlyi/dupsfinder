package main

import (
	"fmt"
	"github.com/svetlyi/dupsfinder/dups"
	"github.com/svetlyi/dupsfinder/file"
	"log"
	"os"
	"sync"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatal("You should have provided path")
		return
	}
	var path string = os.Args[1]
	var procNum int8 = 4
	var filesChannel chan string = make(chan string)
	var filesInfoChannel chan dups.FileInfo = make(chan dups.FileInfo)
	var doneChannel chan bool = make(chan bool)

	go file.WalkThroughFiles(path, &filesChannel)

	wg := sync.WaitGroup{}
	var wgIndex int8
	for wgIndex = 1; wgIndex <= procNum; wgIndex++ {
		wg.Add(1)
		go file.ListenFilesChannel(&filesChannel, &filesInfoChannel, &wg)
	}

	go dups.ListenFilesInfoChannel(&filesInfoChannel, &doneChannel)

	wg.Wait()
	fmt.Println("closing files info channel")
	close(filesInfoChannel)

	<-doneChannel
}
