package logger

import (
	"fmt"
	"github.com/svetlyi/dupsfinder/config"
	"github.com/svetlyi/dupsfinder/structs"
	"io"
	"sync"
)

const (
	TypeError = iota
	TypeMessage
	TypeDelimiter
)

type Logger struct {
	messagesBuf io.Writer
	logMutex    sync.Mutex
	logChan     chan structs.Log
}

func New(writer io.Writer) *Logger {
	return &Logger{
		messagesBuf: writer,
		logMutex:    sync.Mutex{},
		logChan:     make(chan structs.Log),
	}
}

func (logger *Logger) ListenToChannel(exitChan *structs.ExitChan) {
	for logObj := range logger.logChan {
		logger.logMutex.Lock()
		if _, err := logger.messagesBuf.Write([]byte(logObj.Message + "\n")); err != nil {
			fmt.Printf("\ncouldn't write to file %s: %s\n", config.LogFile, err.Error())
			close(*exitChan)
		}
		logger.logMutex.Unlock()
	}
}

func (logger *Logger) Err(msg string) {
	logger.logChan <- structs.Log{
		Type:    TypeError,
		Message: msg,
	}
}

func (logger *Logger) Msg(msg string) {
	logger.logChan <- structs.Log{
		Type:    TypeMessage,
		Message: msg,
	}
}

func (logger *Logger) Delimiter() {
	logger.logChan <- structs.Log{
		Type:    TypeDelimiter,
		Message: "===========================",
	}
}
