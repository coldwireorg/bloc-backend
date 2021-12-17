package files

import (
	"bloc/models"
	"bloc/utils/bcrypto"
	"bloc/utils/errors"
	"io/ioutil"
	"os"

	ecies "github.com/ecies/go"
	"github.com/gofiber/fiber/v2"
)

func Download(ctx *fiber.Ctx) error {
	request := struct {
		AccessId string `json:"accessId"`
	}{}

	// Parse JSON request
	err := ctx.BodyParser(&request)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrRequest)
	}

	access, err := models.AccessGet(request.AccessId)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrDatabaseNotFound)
	}

	file, err := models.FileGet(access.FileId)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrDatabaseNotFound)
	}

	// Check if fil exist
	path := os.Getenv("STORAGE_DIR") + "/" + file.Id
	_, err = os.Stat(path)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}

	// Get private key
	pvKey, err := ecies.NewPrivateKeyFromHex(ctx.Cookies("pvkey"))
	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}

	// get decryption key
	decryptKey, err := ecies.Decrypt(pvKey, access.EncryptionKey)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}

	outfile, err := ioutil.TempFile(os.Getenv("STORAGE_DIR"), file.Id+"*")
	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}
	defer os.Remove(outfile.Name())

	err = bcrypto.DecryptFile(path, decryptKey, int(file.Chunk), func(b []byte) {
		outfile.Write(b)
	})

	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}

	ctx.Response().Header.Add("Content-Disposition", "attachment; filename=\""+file.Name+"\"")
	ctx.Response().Header.Add("Content-Type", file.Type)
	return ctx.SendFile(outfile.Name())
}
