package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/PythonAkoto/base_techtest/adapters/output/logs"
	"github.com/PythonAkoto/base_techtest/env"
)

func StartHTTPServer() {
	logs.Logs(1, "Starting HTTP server...")
	logs.Logs(1, "Loading environment variables...")
	// Load environment variables
	err := env.LoadEnv("env/.env")
	if err != nil {
		logs.Logs(3, fmt.Sprintf("failed to load environment variables: %s", err.Error()))
	}

	// initialise HTTP templates

	// static file server for assets like CSS, JS, images

	// define roiutes and handlers
	http.HandleFunc("/", Hello)

	applicationPort := os.Getenv("APP_PORT") // get application port from environment variable
	logs.Logs(1, fmt.Sprintf("Application port set to: %s", applicationPort))

	// if application port missing from env, default to 9000
	if applicationPort == "" {
		logs.Logs(2, "APP_PORT environment variable not set, defaulting to port 8080")
		applicationPort = "9000" // this can be changed to any default port

		err := http.ListenAndServe(fmt.Sprintf(":%s", applicationPort), nil)
		if err != nil {
			logs.Logs(3, fmt.Sprintf("failed to start HTTP server: %s", err.Error()))
		}
	}

	// start HTTP server via environment variable
	logs.Logs(1, "application started successfully on http://localhost:"+applicationPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", applicationPort), nil)
	if err != nil {
		logs.Logs(3, fmt.Sprintf("failed to start HTTP server: %s", err.Error()))
	}
}
