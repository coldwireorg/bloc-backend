package files

import (
	"github.com/gofiber/fiber/v2"
)

func UnshareFile(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":  "SUCCESS",
		"error": "File unshared!",
	})
}
