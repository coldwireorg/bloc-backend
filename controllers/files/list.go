package files

import (
	"bloc/models"
	"bloc/utils/errors"
	"bloc/utils/tokens"

	"github.com/gofiber/fiber/v2"
)

// List files
func List(ctx *fiber.Ctx) error {
	var err error

	token, err := tokens.Parse(ctx.Cookies("token")) // get JWT token
	if err != nil {
		return errors.HandleError(ctx, errors.ErrRequest)
	}

	files, err := models.FileList(token.Username)

	if err != nil {
		return errors.HandleError(ctx, errors.ErrDatabaseNotFound)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": "SUCCESS",
		"data": fiber.Map{
			"files": files,
		},
	})
}
