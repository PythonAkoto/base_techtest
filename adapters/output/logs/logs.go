package logs

import (
	"log"
)

var logChannel = make(chan string)

const (
	info     = "INFO: "
	warn     = "WARN: "
	errorLog = "ERROR: "
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
func Logs(logType int, message string) {
	var loggedMessage string

	// Determine the log type and format the message accordingly
	switch logType {
	case 1:
		loggedMessage = info + message
	case 2:
		loggedMessage = warn + message
	case 3:
		loggedMessage = errorLog + message
	}

	logChannel <- loggedMessage // Send the log message to the channel for processing
}
