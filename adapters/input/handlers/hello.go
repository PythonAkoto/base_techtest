package handlers

import (
	"fmt"
	"net/http"

	"github.com/PythonAkoto/base_techtest/adapters/output/logs"
)

func Hello(w http.ResponseWriter, r *http.Request) {

	// Respond with a simple message
	fmt.Fprintln(w, "Hello, World!")

	logs.Logs(1, "Home page accessed successfully")
}
