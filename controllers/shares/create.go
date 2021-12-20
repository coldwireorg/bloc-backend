package shares

import (
	"bloc/models"
	"bloc/utils/errors"
	"bloc/utils/tokens"

	ecies "github.com/ecies/go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Create(ctx *fiber.Ctx) error {
	/* 1. Get all needed informations */
	request := struct {
		ShareTo  string `json:"shareTo"`  // Username of the person to share the file to
		FileId   string `json:"fileId"`   // Id of the file to share
		AccessId string `json:"accessId"` // Id of the access of the person who share the file
	}{}

	// Parse JSON request
	err := ctx.BodyParser(&request)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrRequest)
	}

	// Get token
	token, err := tokens.Parse(ctx.Cookies("token")) // get JWT token
	if err != nil {
		return errors.HandleError(ctx, errors.ErrRequest)
	}

	accessExist := models.AccessExist(request.ShareTo, request.FileId)
	println(accessExist)
	if accessExist {
		return errors.HandleError(ctx, errors.ErrDatabaseAlreadyExist)
	}

	/* 3. User verification */
	access, err := models.AccessGet(request.AccessId)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrDatabaseNotFound)
	}

	if access.SharedBy != token.Username {
		return errors.HandleError(ctx, errors.ErrPermission)
	}

	/* 4. File encryption key decryption */
	sharerPrivateKey, err := ecies.NewPrivateKeyFromHex(ctx.Cookies("pvkey"))
	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}

	// get decryption key
	fileKey, err := ecies.Decrypt(sharerPrivateKey, access.EncryptionKey)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}

	/* 5. Re-encrypt file key with the public key of the user to share the file to */
	// Get receiver's public key
	receiverPublicKeyEncoded, err := models.UserGetPubKey(request.ShareTo)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}

	// Get user's public key and transform it to string from byte data
	receiverPublicKey, err := ecies.NewPublicKeyFromBytes(receiverPublicKeyEncoded)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}

	// encrypt file key with receiver's public key
	encryptedFileKey, err := ecies.Encrypt(receiverPublicKey, fileKey)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}

	/* 6. Put sharing in database */
	accessModel := models.Access{
		Id:            uuid.New().String(),
		State:         "SHARED",
		SharedBy:      token.Username,
		SharedTo:      request.ShareTo,
		FileId:        request.FileId,
		Favorite:      false,
		EncryptionKey: encryptedFileKey,
	}

	err = models.AccessCreate(accessModel)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrDatabaseCreate)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": "SUCCESS",
		"data": fiber.Map{
			"file": accessModel,
		},
	})
}
