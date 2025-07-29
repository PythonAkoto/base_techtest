package main

import (
	"log"

	"github.com/PythonAkoto/base_techtest/adapters/input/handlers"
	"github.com/PythonAkoto/base_techtest/adapters/output/logs"
)

func main() {
	// Set log flags to 0 to disable timestamps and other formatting
	log.SetFlags(0)

	// go env.LoadEnv("env/.env") // Load environment variables in a separate goroutine
	go logs.ProcessLogs() // Start processing logs in a separate goroutine

	// go env.LoadEnv(".env")
	go func() {
		handlers.StartHTTPServer() // Start the HTTP server
	}()

	select {} // Block forever to keep the main goroutine alive
}
