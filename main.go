package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)

	fmt.Println("Listening to port 8080!")
	http.ListenAndServe(":8080", mux)
}

func handleRoot(w http.ResponseWriter,
	r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}
