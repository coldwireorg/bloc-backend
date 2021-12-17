package files

import (
	"bloc/models"
	"bloc/utils/errors"
	"bloc/utils/tokens"

	"github.com/gofiber/fiber/v2"
)

// List files
func List(ctx *fiber.Ctx) error {
	var files []*models.FileList
	var filesShared []*models.FileSharedList
	var err error

	shared := ctx.Query("shared")

	token, err := tokens.Parse(ctx.Cookies("token")) // get JWT token
	if err != nil {
		return errors.HandleError(ctx, errors.ErrRequest)
	}

	if shared == "shared" {
		filesShared, err = models.FileListSharedBy(token.Username)
	} else if shared == "received" {
		files, err = models.FileListAll(token.Username, true)
	} else {
		files, err = models.FileListAll(token.Username, false)
	}

	if err != nil {
		return errors.HandleError(ctx, errors.ErrDatabaseNotFound)
	}

	// return the array of file
	if len(files) > 0 {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"code": "SUCCESS",
			"data": files,
		})
	} else {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"code": "SUCCESS",
			"data": filesShared,
		})
	}
}
