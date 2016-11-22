package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	// this matches all paths
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Got request for %q", html.EscapeString(request.URL.Path))
}
