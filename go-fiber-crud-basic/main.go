package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var products []Product

func main() {
	app := fiber.New()

	products = append(products, Product{ID: "1", Name: "Iphone 11", Price: 100000})
	products = append(products, Product{ID: "2", Name: "Iphone 12", Price: 200000})
	products = append(products, Product{ID: "3", Name: "Iphone 13", Price: 300000})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	app.Get("/products", func(c *fiber.Ctx) error {
		return c.JSON(products)
	})

	app.Get("/products/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for _, product := range products {
			if product.ID == id {
				return c.JSON(product)
			}
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	})

	app.Post("/products", func(c *fiber.Ctx) error {
		var product Product
		if err := c.BodyParser(&product); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}
		products = append(products, product)
		return c.JSON(product)
	})
	
	if err := app.Listen(":5655"); err != nil {
		panic(err)
	}

	func getallProducts(c *fiber.Ctx) error {
		return c.JSON(products)
	}

	funcgetProductByID(c *fiber.Ctx) error {
		id := c.Params("id")
		for _, product := range products {
			if product.ID == id {
				return c.JSON(product)
			}

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	func createProduct(c *fiber.Ctx) error {
		var product Product
		if err := c.BodyParser(&product); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}
		products = append(products, product)
		return c.JSON(product)
	}

	func updateProductByID(c *fiber.Ctx) error {
		id := c.Params("id")

		if err := c.BodyParser(&product); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		product := new(Product)

		if err:= c.BodyParser(product); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		for i, producData := range products {
			if producData.ID == id {
				products[i] = *product
				return c.JSON(product)
			}
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}
}
}