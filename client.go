package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	print("test")


	rand.Seed(time.Now().UnixNano())

	producten := []string{"kaas", "apple", "banana", "cherry", "date", "elderberry"}


	fmt.Println(randomItem)

	while (true) {
		randomIndex := rand.Intn(len(producten))
		randomItem := producten[randomIndex]

		// get product aantal

		// koop wanneer voorradig

		time.Sleep(1 * time.Second)
	}

}
