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

func main() {
	rand.Seed(time.Now().UnixNano())

	producten := []string{"kaas", "apple", "banana", "cherry", "date", "elderberry"}

	for {

		randomIndex := rand.Intn(len(producten)) // Get a random index from the array
		randomItem := producten[randomIndex]
		fmt.Println("Buying item:", randomItem)

		// get product aantal
		product := randomItem // replace with your actual product
		url := fmt.Sprintf("http://localhost:8000/product/%s", product)

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
		//fmt.Println(string(body))

		aantal, err := strconv.Atoi(string(body))

		fmt.Println("Stock item:", aantal)

		// koop wanneer voorradig

		if aantal < 1 {
			fmt.Println("Error: Out of stock")
			time.Sleep(1 * time.Second)

			continue
		} else {

			// Construct the URL with the product and specify the port
			url := fmt.Sprintf("http://localhost:8000/product/%s", product)

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
