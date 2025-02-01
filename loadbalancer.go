package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Received request: %s : %s", r.Method, r.URL.Path)
}
func main() {
	fmt.Println("Starting load balancer")
	port := "80"
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("error starting server:", err)
	}
}
