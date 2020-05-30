package main

import (
	"log"
	"net/http"

	app "github.com/hirakiuc/gonta-app"
)

func main() {
	http.HandleFunc("/serve", app.Serve)

	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal(err)
	}
}
