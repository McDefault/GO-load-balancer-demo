package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var mu sync.Mutex
var producten = map[string]int{
	"kaas":   10,
	"worst":  2,
	"apple":  3,
	"banana": 1,
	"cherry": 0,
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	paramValue := strings.TrimPrefix(r.URL.Path, "/product/")
	switch r.Method {
	case "GET":

		fmt.Println("Buying item:", paramValue, producten[paramValue])

		fmt.Fprintf(w, strconv.Itoa(producten[paramValue]))
	case "POST":
		mu.Lock()
		defer mu.Unlock()

		if producten[paramValue] < 0 {
			fmt.Println("An error occurred: voorraad onder 0")
			os.Exit(1) // Exit with status code 1
		}
		if producten[paramValue] < 1 {
			fmt.Fprintf(w, "Product: %s is OP! %d", paramValue, producten[paramValue])
			fmt.Println("Error: Out of stock")
		} else {
			producten[paramValue]--
			fmt.Fprintf(w, "Item bought: %s; %d, left", paramValue, producten[paramValue])
		}
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/product/", productHandler)

	fmt.Println("Server starting at http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
