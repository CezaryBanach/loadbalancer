package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type Server struct {
	URL    string
	Health bool
	mu     sync.Mutex
}

func (s *Server) GetHealth() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Health
}

func NewServer(url string) *Server {
	return &Server{
		URL:    url,
		Health: true,
	}
}

var servers = []*Server{
	NewServer("http://localhost:8081"),
	NewServer("http://localhost:8082"),
	NewServer("http://localhost:8083"),
}

var cur_serv = 0

func get_next_server() *Server {
	response := servers[cur_serv]
	cur_serv = (cur_serv + 1) % len(servers)
	return response
}

func get_backend() (*url.URL, error) {
	first_check := cur_serv
	server := get_next_server()
	for server.Health == false {
		server = get_next_server()
		//we have tried all servers and are back to first one
		if cur_serv == first_check {
			return nil, fmt.Errorf("All servers unhealthy")
		}
	}

	url, err := url.Parse(server.URL)
	if err != nil {
		return nil, err
	}
	return url, nil
}
func print_request(r *http.Request) {
	fmt.Printf("\n##################")
	fmt.Printf("Received request from %s\n%s %s %s\nHost: %s\nUser-Agent: %s\n", r.RemoteAddr, r.Method, r.URL.Path, r.Proto, r.Host, r.UserAgent())
}
func respond(w http.ResponseWriter, resp *http.Response, body string) {
	fmt.Printf("Response from server: %s %d\n%s\n", resp.Proto, resp.StatusCode, body)
	w.WriteHeader(resp.StatusCode)
	fmt.Fprintf(w, string(body))
}
func proxy_request(r *http.Request, url *url.URL) *http.Request {
	r.URL = url
	//can't have this set
	r.RequestURI = ""
	return r
}

func handle_request(w http.ResponseWriter, r *http.Request) {
	print_request(r)
	backend, err := get_backend()
	if err != nil {
		fmt.Println("error getting backend", err)
		//return error etc
		return
	}
	fmt.Println("sending to server", backend.String())

	client := http.DefaultClient
	time.Sleep(time.Second * 10)

	r = proxy_request(r, backend)
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	//send the response to the caller
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response", err)
		return
	}

	respond(w, resp, string(body))
}

func handler(w http.ResponseWriter, r *http.Request) {
	handle_request(w, r)
	fmt.Println("Finished processing request")
	fmt.Println("Exiting http handler")
}

func health_check(server *Server) bool {
	resp, err := http.Get(server.URL + "/healthy")
	if err != nil {
		fmt.Printf("Health check failed for %s\n", server.URL)
		return false
	}
	defer resp.Body.Close()
	fmt.Println("Got status code", resp.StatusCode, server.URL)
	if resp.StatusCode != 200 {
		return false
	}
	return true
}
func start_health_checks() {
	for _, server := range servers {
		cur_server := server
		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				healthy := health_check(cur_server)
				cur_server.mu.Lock()
				cur_server.Health = healthy
				cur_server.mu.Unlock()
			}
		}()
	}
}

func main() {
	fmt.Println("Starting load balancer")
	port := "80"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	if len(os.Args) > 2 {
		port = os.Args[1]
	}
	fmt.Println("starting health checks")
	start_health_checks()
	fmt.Printf("Listening on port %s\n", port)
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("error starting server:", err)
	}
}
