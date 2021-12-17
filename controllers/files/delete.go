package files

import (
	"bloc/models"
	"bloc/utils"
	"bloc/utils/errors"
	"bloc/utils/tokens"
	"os"

	"github.com/gofiber/fiber/v2"
)

// Delete file
func Delete(ctx *fiber.Ctx) error {
	var total int64

	request := struct {
		FileId   string `json:"fileId"`
		AccessId string `json:"accessId"`
	}{}

	// Parse JSON request
	err := ctx.BodyParser(&request)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrRequest)
	}

	// parse user's token
	token, err := tokens.Parse(ctx.Cookies("token"))
	if err != nil {
		return errors.HandleError(ctx, errors.ErrRequest)
	}

	// Get user quota to update its quota
	userQuota, err := models.UserGetQuota(token.Username)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrDatabaseNotFound)
	}

	file, err := models.FileGet(request.FileId)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrDatabaseNotFound)
	}

	err = models.AccessDelete(request.AccessId, token.Username)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrDatabaseRemove)
	}

	// if the person if the owner of the file, delete the file
	if file.Owner == token.Username {
		err = models.FileDelete(request.FileId, token.Username)
		if err != nil {
			return errors.HandleError(ctx, errors.ErrDatabaseRemove)
		}

		// remove file on the user side and update the their quota.
		total = userQuota - file.Size

		if total < 0 {
			total = 0
		}

		err = models.UserUpdateQuota(token.Username, total)
		if err != nil {
			return errors.HandleError(ctx, errors.ErrDatabaseUpdate)
		}

		path := os.Getenv("STORAGE_DIR") + "/" + request.FileId
		err = os.Remove(path)
		if err != nil {
			return errors.HandleError(ctx, errors.ErrInternal)
		}
	}

	return ctx.JSON(fiber.Map{
		"code":    "SUCCESS",
		"message": "File deleted",
		"quota": fiber.Map{
			"max":   utils.GetQuota(),
			"total": total,
		},
	})
}
