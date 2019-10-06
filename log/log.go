package log

import "sync"

const (
	TypeError     = 1
	TypeMessage   = 2
	TypeDelimiter = 3
)

type Log struct {
	Type    int8
	Message string
}

var messages []string
var logMutex = sync.Mutex{}

func ListenToChannel(logChannel *chan Log) {
	for logObj := range *logChannel {
		logMutex.Lock()
		messages = append(messages, logObj.Message)
		logMutex.Unlock()
	}
}

func GetLastMessages(num int) []string {
	if len(messages) < num {
		return messages
	}
	return messages[len(messages)-num:]
}

func Err(logChannel *chan Log, msg string) {
	*logChannel <- Log{
		Type:    TypeError,
		Message: msg,
	}
}

func Msg(logChannel *chan Log, msg string) {
	*logChannel <- Log{
		Type:    TypeMessage,
		Message: msg,
	}
}

func Delimiter(logChannel *chan Log) {
	*logChannel <- Log{
		Type:    TypeDelimiter,
		Message: "===========================",
	}
}
