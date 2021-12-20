package shares

import "github.com/gofiber/fiber/v2"

func List(ctx *fiber.Ctx) error {
	return ctx.SendStatus(404)
}
