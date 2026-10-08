package main

import ("github.com/gofiber/fiber/v3")

func main() {
	app := fiber.New()

	app.Get("/hellogo", func(c fiber.Ctx) error {
		return c.SendString("HelloGo....")
	})

	app.Get("/higo", func(c fiber.Ctx) error {
		return c.SendString("HiGo....")
	})

	if err :=app.Listen(":5599"); err != nil {
		panic(err)
	}
}