package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	loadBalancerURL := flag.String("lb", "http://localhost:8000", "Load balancer URL")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	for {

		// Measure the start time
		startTime := time.Now()

		// Randomly select a number between 0 and 1000
		randomNumber := rand.Intn(1001)
		fmt.Println("Sending number:", randomNumber)

		// Construct the URL to send the request to the load balancer
		url := fmt.Sprintf(`%s/%d`, *loadBalancerURL, randomNumber)

		// Make a GET request
		resp, err := http.Get(url)

		// Measure the end time
		endTime := time.Now()
		responseTime := endTime.Sub(startTime)

		// If there is an error with the request, print it to the console
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Defer the closing of the body
		defer resp.Body.Close()

		// Calculate the performance metric
		performanceMetric := responseTime.Milliseconds() - int64(randomNumber)

		fmt.Printf("Response time: %v ms\n", performanceMetric)
		//time.Sleep(1 * time.Second)

	}
}
