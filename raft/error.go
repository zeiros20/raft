// raft/error.go
package raft

import (
	"log"
	"os"
)

var (
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	InfoLogger  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
)

// HandleError prints the error if non-nil and returns whether it was handled
func HandleError(err error, context string) bool {
	if err != nil {
		ErrorLogger.Printf("%s: %v", context, err)
		return true
	}
	return false
}

// FatalError logs the error and exits
func FatalError(err error, context string) {
	if err != nil {
		ErrorLogger.Fatalf("%s: %v", context, err)
	}
}

// InfoLogger logs informational messages
func InfoMessage(message string) {
	InfoLogger.Println(message)
}
