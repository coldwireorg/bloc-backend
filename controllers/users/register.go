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

	"github.com/alexedwards/argon2id"
	ecies "github.com/ecies/go"
	"github.com/gofiber/fiber/v2"
)

func Register(ctx *fiber.Ctx) error {
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

	// Verify that the username is not empty or too short
	if request.Username == "" {
		return errors.HandleError(ctx, errors.ErrRequest)
	}

	// Verify if the user already exist
	exist := models.UserExist(request.Username)
	log.Println(exist)
	if exist {
		return errors.HandleError(ctx, errors.ErrAuthExist)
	}

	// Hash password with argon2id
	hash, err := argon2id.CreateHash(request.Password, argon2id.DefaultParams)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}

	// Generate user's keypair
	pvKey, err := ecies.GenerateKey()
	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}

	// Encrypt private key with user's password wich is derived using Argon2i
	pvKeyEncKey := bcrypto.DeriveKey([]byte(request.Password))
	pvKeyEncrypted, err := bcrypto.Encrypt(pvKey.Bytes(), pvKeyEncKey)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}

	err = models.UserCreate(request.Username, hash, pvKey.PublicKey.Bytes(false), pvKeyEncrypted)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrDatabaseCreate)
	}

	exp := time.Hour * 2                          // define token expiration
	jwt := tokens.Generate(request.Username, exp) // generate JWT token

	// Set cookies
	ctx.Cookie(utils.GenCookie("token", jwt, exp, os.Getenv("SERVER_DOMAIN")))
	ctx.Cookie(utils.GenCookie("pvkey", pvKey.Hex(), exp, os.Getenv("SERVER_DOMAIN")))
	// If the user use TOR, set the cookies on the tor address
	if ctx.Hostname() == os.Getenv("TOR_ADDRESS") {
		ctx.Cookie(utils.GenCookie("token", jwt, exp, os.Getenv("TOR_ADDRESS")))
		ctx.Cookie(utils.GenCookie("pvkey", pvKey.Hex(), exp, os.Getenv("TOR_ADDRESS")))
	}

	// Return login informations
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":    "SUCCESS",
		"message": "User created",
		"data": fiber.Map{
			"username": request.Username,
			"token":    jwt,
			"quota": fiber.Map{
				"max":   utils.GetQuota(),
				"total": 0,
			},
		},
	})
}
