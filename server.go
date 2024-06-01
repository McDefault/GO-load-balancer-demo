package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"strconv"
)

func main() {
	print("test")

	//producten := make(map[string]int)
	//
	//producten["kaas"] = 7
	//producten["worst"] = 13

	producten := map[string]int{
		"kaas":   5,
		"worst":  15,
		"apple":  5,
		"banana": 7,
		"cherry": 15,
	}

	fmt.Println("App running")

	app := fiber.New()
	app.Use(cors.New())

	api := app.Group("/api")

	// Test handler
	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("App running")
	})

	api.Get("/:param", func(c *fiber.Ctx) error {
		paramValue := c.Params("param")
		return c.SendString(strconv.Itoa(producten[paramValue]))
	})

	api.Post("/:param", func(c *fiber.Ctx) error {

		paramValue := c.Params("param")
		voorraad := producten[paramValue]
		if voorraad < 1 {
			return c.SendString("Product: " + paramValue + " is OP! " + strconv.Itoa(producten[paramValue]))
		}
		producten[paramValue]--
		return c.SendString("Parameter value: " + paramValue + " = " + strconv.Itoa(producten[paramValue]))
	})

	log.Fatal(app.Listen(":5000"))

}
