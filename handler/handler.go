package handler

import (
	"fmt"
	"net/http"
)

// ServeHTTP handles the http request.
func Serve(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
