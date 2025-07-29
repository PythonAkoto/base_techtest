package logs

import (
	"fmt"
	"log"
	"time"
)

var logChannel = make(chan string)

const (
	info     = "INFO"
	warn     = "WARN"
	errorLog = "ERROR"
)

// ProcessLogs reads log messages from the logChannel and prints them using the standard logger.
// It runs indefinitely, processing messages as they are received.
func ProcessLogs() {
	for logMessage := range logChannel {
		log.Println(logMessage)
	}
}

/*
Logs takes a log type and a message and sends a formatted log message to the logging channel.

The log type determines the type of log message to be sent. Currently supported types are:

1: INFO

2: WARN

3: ERROR

The message should be the log message text. The function will add the appropriate prefix and send the message to the logging channel.
*/
func Logs(logType int, message string, provider string) {
	var loggedMessage string
	// Use custom timestamp format
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Dynamically determine the log type and format the message accordingly
	if provider != "" {
		switch logType {
		case 1:
			loggedMessage = fmt.Sprintf("\"level\":%q, \"provider\":%q, \"timestamp\":%q, \"message\":%q", info, provider, timestamp, message)
		case 2:
			loggedMessage = fmt.Sprintf("\"level\":%q, \"provider\":%q, \"timestamp\":%q, \"message\":%q", warn, provider, timestamp, message)
		case 3:
			loggedMessage = fmt.Sprintf("\"level\":%q, \"provider\":%q, \"timestamp\":%q, \"message\":%q", errorLog, provider, timestamp, message)
		}
	} else {
		switch logType {
		case 1:
			loggedMessage = fmt.Sprintf("\"level\":%q, \"timestamp\":%q, \"message\":%q", info, timestamp, message)
		case 2:
			loggedMessage = fmt.Sprintf("\"level\":%q, \"timestamp\":%q, \"message\":%q", warn, timestamp, message)
		case 3:
			loggedMessage = fmt.Sprintf("\"level\":%q, \"timestamp\":%q, \"message\":%q", errorLog, timestamp, message)
		}
	}

	logChannel <- loggedMessage // Send the log message to the channel for processing
}
