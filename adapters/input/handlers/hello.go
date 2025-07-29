package handlers

import (
	"fmt"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {

	// Respond with a simple message
	fmt.Fprintln(w, "Hello, World!")
}
