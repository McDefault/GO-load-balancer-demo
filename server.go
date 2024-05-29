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

	m := make(map[string]int)

	m["kaas"] = 7
	m["worst"] = 13

	fmt.Println("App running")

	app := fiber.New()
	app.Use(cors.New())

	api := app.Group("/api")

	// Test handler
	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("App running")
	})

	//api.Get("/:param", func(c *fiber.Ctx) error {
	//	paramValue := c.Params("param")
	//	return c.SendString("Parameter value: " + paramValue + " = " + strconv.Itoa(m[paramValue]))
	//})

	api.Post("/:param", func(c *fiber.Ctx) error {

		paramValue := c.Params("param")
		voorraad := m[paramValue]
		if voorraad < 1 {
			return c.SendString("Product: " + paramValue + " is OP!" + strconv.Itoa(m[paramValue]))
		}
		m[paramValue]--
		return c.SendString("Parameter value: " + paramValue + " = " + strconv.Itoa(m[paramValue]))
	})

	log.Fatal(app.Listen(":5000"))

}
