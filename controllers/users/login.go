package users

import (
	"bloc/models"
	"bloc/utils"
	"bloc/utils/bcrypto"
	"bloc/utils/errors"
	"bloc/utils/tokens"
	"log"
	"os"
	"time"

	ecies "github.com/ecies/go"

	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
)

func Login(ctx *fiber.Ctx) error {
	// Structure of the JSON request
	request := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	// Parse JSON request
	err := ctx.BodyParser(&request)
	if err != nil {
		log.Println(err, request)
		return errors.HandleError(ctx, errors.ErrRequest)
	}

	// Verify that the username is not empty
	if request.Username == "" {
		return errors.HandleError(ctx, errors.ErrRequest)
	}

	// Get the user
	user, err := models.UserGet(request.Username)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrDatabaseNotFound)
	}

	// Verify password
	isValidPassword, err := argon2id.ComparePasswordAndHash(request.Password, user.Password)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}

	if !isValidPassword {
		return errors.HandleError(ctx, errors.ErrAuthPassword)
	}

	// Decrypt private key with password
	decryptedPrivateKey, err := bcrypto.Decrypt(user.PrivateKey, bcrypto.DeriveKey([]byte(request.Password)))
	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}

	// Get private key
	privateKey := ecies.NewPrivateKeyFromBytes(decryptedPrivateKey)

	exp := time.Hour * 2
	jwt := tokens.Generate(request.Username, privateKey.Hex(), exp) // Generate JWT token

	// set cookies
	ctx.Cookie(utils.GenCookie("token", jwt, exp, os.Getenv("SERVER_DOMAIN")))
	// if user use tor browser
	if ctx.Hostname() == os.Getenv("TOR_ADDRESS") {
		ctx.Cookie(utils.GenCookie("token", jwt, exp, os.Getenv("TOR_ADDRESS")))
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":    "SUCCESS",
		"message": "User login",
		"data": fiber.Map{
			"username": request.Username,
			"quota": fiber.Map{
				"total": 0,
				"max":   utils.GetQuota(),
			},
		},
	})
}
