package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/PythonAkoto/base_techtest/adapters/output/logs"
	"github.com/PythonAkoto/base_techtest/env"
)

func StartHTTPServer() {
	logs.Logs(1, "Starting HTTP server...", "")
	logs.Logs(1, "Loading environment variables...", "")
	// Load environment variables
	err := env.LoadEnv("env/.env")
	if err != nil {
		logs.Logs(3, fmt.Sprintf("failed to load environment variables: %s", err.Error()), "")
	}

	// initialise HTTP templates

	// static file server for assets like CSS, JS, images

	// define roiutes and handlers
	http.HandleFunc("/", Hello)
	http.HandleFunc("/products", GetProductsHandler)

	applicationPort := os.Getenv("APP_PORT") // get application port from environment variable
	
	// if application port missing from env, default to 8080
	if applicationPort == "" {
		logs.Logs(2, "APP_PORT environment variable not set, defaulting to port 8080", "")
		applicationPort = "8080" // default port
	}
	
	logs.Logs(1, fmt.Sprintf("Application port set to: %s", applicationPort), "")

	// start HTTP server
	logs.Logs(1, "application started successfully on http://localhost:"+applicationPort, "")
	err = http.ListenAndServe(fmt.Sprintf(":%s", applicationPort), nil)
	if err != nil {
		logs.Logs(3, fmt.Sprintf("failed to start HTTP server: %s", err.Error()), "")
	}
}
