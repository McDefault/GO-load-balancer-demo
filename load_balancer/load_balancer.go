package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server struct {
	URL         *url.URL
	IsHealthy   bool
	Mutex       sync.Mutex
	Connections int
}

type LoadBalancer struct {
	Servers  []*Server
	Strategy LoadBalancingStrategy
	Mutex    sync.Mutex
}

type LoadBalancingStrategy interface {
	GetNextServer(servers []*Server) *Server
}

func NewLoadBalancer(servers []*Server, strategy LoadBalancingStrategy) *LoadBalancer {
	return &LoadBalancer{
		Servers:  servers,
		Strategy: strategy,
	}
}

// load balancers

type RoundRobinStrategy struct {
	current int
	mutex   sync.Mutex
}

func (rr *RoundRobinStrategy) GetNextServer(servers []*Server) *Server {
	rr.mutex.Lock()
	defer rr.mutex.Unlock()
	server := servers[rr.current]

	rr.current = (rr.current + 1) % len(servers)
	return server
}

type WeightedRoundRobinStrategy struct {
	weights []int
	current int
	mutex   sync.Mutex
}

func (wrr *WeightedRoundRobinStrategy) GetNextServer(servers []*Server) *Server {
	wrr.mutex.Lock()
	defer wrr.mutex.Unlock()
	totalWeight := 0
	for _, weight := range wrr.weights {
		totalWeight += weight
	}
	wrr.current = (wrr.current + 1) % totalWeight
	for i, weight := range wrr.weights {
		if wrr.current < weight {
			return servers[i]
		}
		wrr.current -= weight
	}
	return servers[0]
}

type LeastConnectionStrategy struct {
	mutex sync.Mutex
}

func (lc *LeastConnectionStrategy) GetNextServer(servers []*Server) *Server {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	var leastConnServer *Server
	minConnections := int(^uint(0) >> 1) // Initialize to max int value

	for _, server := range servers {
		server.Mutex.Lock()
		if server.IsHealthy && server.Connections < minConnections {
			leastConnServer = server
			minConnections = server.Connections
		}
		server.Mutex.Unlock()
	}

	return leastConnServer
}

//
//func (lb *LoadBalancer) GetNextServer() *Server {
//	lb.Mutex.Lock()
//	defer lb.Mutex.Unlock()
//
//	server := lb.Servers[lb.Current]
//	lb.Current = (lb.Current + 1) % len(lb.Servers)
//	return server
//}

//func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	server := lb.GetNextServer()
//	proxy := httputil.NewSingleHostReverseProxy(server.URL)
//	proxy.ServeHTTP(w, r)
//}

func parseURL(rawURL string) *url.URL {
	url, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return url
}

func (lb *LoadBalancer) AddServer(server *Server) {
	lb.Mutex.Lock()
	defer lb.Mutex.Unlock()
	lb.Servers = append(lb.Servers, server)
	fmt.Println(`Adding server `, *server.URL)

}

func (lb *LoadBalancer) EnlistServerHandler(w http.ResponseWriter, r *http.Request) {
	rawURL := r.URL.Query().Get("url")
	if rawURL == "" {
		http.Error(w, "URL parameter is required", http.StatusBadRequest)
		return
	}

	serverURL, err := url.Parse(rawURL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	server := &Server{URL: serverURL, IsHealthy: true}
	lb.AddServer(server)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server added successfully"))
}

func (lb *LoadBalancer) SetStrategy(strategy LoadBalancingStrategy) {
	lb.Mutex.Lock()
	defer lb.Mutex.Unlock()
	lb.Strategy = strategy
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := lb.Strategy.GetNextServer(lb.Servers)

	server.Mutex.Lock()
	server.Connections++
	server.Mutex.Unlock()
	fmt.Printf("Server: %v, %v\n", server.URL, server.Connections)

	proxy := httputil.NewSingleHostReverseProxy(server.URL)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		fmt.Printf("Error proxying request to server %s: %v\n", server.URL, err)
		http.Error(w, "Error forwarding request", http.StatusBadGateway)

		server.Mutex.Lock()
		server.Connections--
		server.Mutex.Unlock()

	}

	proxy.ServeHTTP(w, r)

	server.Mutex.Lock()
	server.Connections--
	server.Mutex.Unlock()
}

func main() {
	port := flag.String("port", "8000", "Port to run the server on")
	flag.Parse()

	servers := []*Server{
		{URL: parseURL("http://localhost:8080"), IsHealthy: true},
	}

	lb := NewLoadBalancer(servers, &LeastConnectionStrategy{})

	http.HandleFunc("/enlist", lb.EnlistServerHandler)
	http.Handle("/", lb)

	fmt.Printf("Starting server on port %s\n", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}

	//fmt.Println("Server starting at http://localhost:8080")
	//http.ListenAndServe(":8080", nil);

}
