package controllers

import "github.com/gofiber/fiber/v2"

func TestUserControler(c *fiber.Ctx) error {
	return c.SendString("TEST USERC ONTROLLER")
}
