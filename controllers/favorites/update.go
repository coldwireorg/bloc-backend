package favorites

import (
	"bloc/models"
	"bloc/utils/errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Update(ctx *fiber.Ctx) error {
	request := struct {
		AccessId string `json:"accessId"`
		Favorite bool   `json:"favorite"`
	}{}

	// Parse JSON request
	err := ctx.BodyParser(&request)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrRequest)
	}

	err = models.FileUpdateFavorite(request.Favorite, request.AccessId)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}

	return ctx.JSON(fiber.Map{
		"code":    "SUCCESS",
		"message": fmt.Sprint("favorite set to:", request.Favorite),
	})
}
