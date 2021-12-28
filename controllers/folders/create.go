package folders

import (
	"bloc/models"
	"bloc/utils/errors"
	"bloc/utils/tokens"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var request struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func Create(ctx *fiber.Ctx) error {
	// Parse JSON request
	err := ctx.BodyParser(&request)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrRequest)
	}

	token, err := tokens.Parse(ctx.Cookies("token")) // Parse user's JWT token
	if err != nil {
		return errors.HandleError(ctx, errors.ErrRequest)
	}

	if request.Name == "" || request.Path == "" {
		return errors.HandleError(ctx, errors.ErrRequest)
	}

	folder := models.Folder{
		Id:    uuid.New().String(),
		Owner: token.Username,
		Name:  request.Name,
		Path:  request.Path,
	}

	err = models.FolderCreate(folder)
	if err != nil {
		println(err)
		return errors.HandleError(ctx, errors.ErrInternal)
	}

	return ctx.JSON(fiber.Map{
		"code":    "SUCCESS",
		"message": "File deleted",
		"data": fiber.Map{
			"folder": folder,
		},
	})
}
