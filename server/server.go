package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	port := flag.String("port", "8080", "Port to run the server on")
	loadBalancerURL := flag.String("lb", "http://localhost:8000/enlist", "Load balancer enlist URL")
	flag.Parse()

	// Enlist the server to the load balancer
	enlistServer(*loadBalancerURL, *port)

	http.HandleFunc("/", numberHandler)
	fmt.Printf("Starting server on port %s\n", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}

}

func numberHandler(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/")
	number, err := strconv.Atoi(path)
	if err != nil || number < 0 || number > 1000 {
		http.Error(w, "Invalid number. It must be a number between 0 and 1000.", http.StatusBadRequest)
		return
	}
	duration := time.Duration(number) * time.Millisecond
	fmt.Printf("Waiting %v ms\n", duration)
	time.Sleep(duration)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Status: 200 OK\n")

}

func enlistServer(loadBalancerURL, port string) {
	url := fmt.Sprintf("%s?url=http://localhost:%s", loadBalancerURL, port)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error enlisting server: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to enlist server: %s\n", resp.Status)
	} else {
		fmt.Println("Server enlisted successfully")
	}
}
