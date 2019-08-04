package main

import (
	"net/http"

	app "github.com/hirakiuc/gonta-app"
)

func main() {
	http.HandleFunc("/serve", app.Serve)
	http.ListenAndServe(":8082", nil)
}
