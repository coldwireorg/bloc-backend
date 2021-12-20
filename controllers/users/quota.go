package users

import (
	"bloc/models"
	"bloc/utils"
	"bloc/utils/errors"
	"bloc/utils/tokens"

	"github.com/gofiber/fiber/v2"
)

// List files
func QuotaCheck(ctx *fiber.Ctx) error {
	request := struct {
		FileSize int64 `json:"size"`
	}{}

	// Parse JSON request
	err := ctx.BodyParser(&request)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrRequest)
	}

	token, err := tokens.Parse(ctx.Cookies("token")) // get JWT token
	if err != nil {
		return errors.HandleError(ctx, errors.ErrRequest)
	}

	userQuota, err := models.UserGetQuota(token.Username)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrDatabaseNotFound)
	}

	if userQuota+request.FileSize > int64(utils.GetQuota()) {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "ERROR_QUOTA",
			"message": "You don't have enough space",
			"quota": fiber.Map{
				"max":   utils.GetQuota(),
				"total": userQuota,
			},
		})
	}

	// tell the frontend it's possible to upload the file
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    "SUCCESS",
		"message": "you can upload your file",
		"data": fiber.Map{
			"quota": fiber.Map{
				"max":   utils.GetQuota(),
				"total": userQuota,
			},
		},
	})
}

func QuotaGet(ctx *fiber.Ctx) error {
	token, err := tokens.Parse(ctx.Cookies("token")) // get JWT token
	if err != nil {
		return errors.HandleError(ctx, errors.ErrRequest)
	}

	userQuota, err := models.UserGetQuota(token.Username)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrDatabaseNotFound)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": "SUCCESS",
		"data": fiber.Map{
			"quota": fiber.Map{
				"max":   utils.GetQuota(),
				"total": userQuota,
			},
		},
	})
}
