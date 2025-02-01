package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

var port = "80"

func get_response() (string, error) {
	content, err := os.ReadFile("backend/response_template.txt")
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(content), "$PORT", port), nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received request from %s\n%s %s %s\nHost: %s\nUser-Agent: %s\n", r.RemoteAddr, r.Method, r.URL.Path, r.Proto, r.Host, r.UserAgent())
	w.WriteHeader(200)
	response, err := get_response()
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
	}
	fmt.Fprintf(w, response)
	fmt.Printf("responded\n")
}
func main() {
	fmt.Println("Starting backend server...")
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	fmt.Printf("Listening on port %s\n", port)
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("error starting server:", err)
	}
}
