package dups

import "fmt"

func ListenFilesInfoChannel(filesInfoChannel *chan FileInfo, doneChannel *chan bool) {
	var fileDups = make(map[string][]FileInfo)

	for fileInfo := range *filesInfoChannel {
		if nil == fileDups[fileInfo.Hash] {
			fileDups[fileInfo.Hash] = make([]FileInfo, 0)
		}
		fileDups[fileInfo.Hash] = append(fileDups[fileInfo.Hash], fileInfo)
	}

	printDups(fileDups)
	*doneChannel <- true
}

func printDups(files map[string][]FileInfo){
	fmt.Println("==================================")
	fmt.Println("Dups:")

	for _, sameHashFiles := range files {
		if len(sameHashFiles) > 1 {
			fmt.Printf("Found dups for %v: \n", sameHashFiles[0].Path)

			for _, file := range sameHashFiles {
				fmt.Println(file.Path)
			}
		}
	}
}