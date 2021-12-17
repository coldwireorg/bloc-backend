package errors

import "github.com/gofiber/fiber/v2"

type apiError struct {
	status  int
	code    string
	message string
}

func HandleError(ctx *fiber.Ctx, err apiError) error {
	return ctx.Status(err.status).JSON(fiber.Map{
		"code":    err.code,
		"message": err.message,
	})
}
