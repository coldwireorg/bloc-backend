package files

import (
	"bloc/models"
	"bloc/utils/errors"
	"bloc/utils/tokens"

	"github.com/gofiber/fiber/v2"
)

// List files
func List(ctx *fiber.Ctx) error {
	var path string

	path = ctx.Query("path") // Get path
	if path == "" {
		path = "/"
	}

	token, err := tokens.Parse(ctx.Cookies("token")) // get JWT token
	if err != nil {
		return errors.HandleError(ctx, errors.ErrRequest)
	}

	files, err := models.FileList(token.Username, path)

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
