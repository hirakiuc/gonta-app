package main

import (
	"net/http"

	"github.com/hirakiuc/gonta-app/handler"
)

func main() {
	http.HandleFunc("/serve", handler.Serve)
	http.ListenAndServe(":8082", nil)
}
