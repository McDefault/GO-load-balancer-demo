package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// getRandomKey is a helper function that returns a random key from the given map.
func getRandomKey(m map[string]int) string {
	// Create a slice to hold the keys from the map.
	keys := make([]string, len(m))

	// Fill the slice with keys from the map.
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	// Return a random key from the slice.
	return keys[rand.Intn(len(keys))]
}

func main() {
	print("test")

	rand.Seed(time.Now().UnixNano())

	producten := []string{"kaas", "apple", "banana", "cherry", "date", "elderberry"}

	for {

		randomIndex := rand.Intn(len(producten)) // Get a random index from the array
		randomItem := producten[randomIndex]
		fmt.Println("Random entry:", randomIndex, randomItem)

		// get product aantal
		product := randomItem // replace with your actual product
		url := fmt.Sprintf("http://localhost:5000/api/%s", product)

		// Make a GET request
		resp, err := http.Get(url)

		// If there is an error with the request, print it to the console
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Defer the closing of the body
		defer resp.Body.Close()

		// Read the content into a byte array
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// At this point, `body` contains the raw server response
		fmt.Println(string(body))

		aantal, err := strconv.Atoi(string(body))

		// koop wanneer voorradig

		if aantal < 1 {
			fmt.Println("Error: Out of stock")

			continue
		} else {

			// Construct the URL with the product and specify the port
			url := fmt.Sprintf("http://localhost:5000/api/%s", product)

			// Define the data for the POST request
			data := []byte(`{"productID": "123", "productName": "apple"}`)

			// Create a new request using http
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))

			// If there is an error with the request, print it to the console
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			// Add headers
			req.Header.Add("Content-Type", "application/json")

			// Send the request via a client
			client := &http.Client{}
			resp, err := client.Do(req)

			// Handle the error if there is one
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			// Defer the closing of the body
			defer resp.Body.Close()

			// Read the content into a byte array
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			// At this point, `body` contains the raw server response
			fmt.Println(string(body))
		}

		time.Sleep(1 * time.Second)

	}

}
