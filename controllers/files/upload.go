package files

import (
	"bloc/models"
	"bloc/utils"
	"bloc/utils/bcrypto"
	"bloc/utils/errors"
	"bloc/utils/tokens"
	"os"
	"time"

	ecies "github.com/ecies/go"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/chacha20poly1305"
)

// File upload function
func Upload(ctx *fiber.Ctx) error {
	fileMultipart, err := ctx.FormFile("file") // We are getting the file sent via a form
	if err != nil {
		// If we can't get the sent file, we sen a 500 error
		return errors.HandleError(ctx, errors.ErrRequest)
	}

	token, err := tokens.Parse(ctx.Cookies("token")) // Parse user's JWT token
	if err != nil {
		return errors.HandleError(ctx, errors.ErrRequest)
	}

	// Get user's data in the database
	user, err := models.UserGet(token.Username)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrDatabaseNotFound)
	}

	// Check if the user have enought space
	quota := user.Quota + fileMultipart.Size
	if quota > int64(utils.GetQuota()) {
		return errors.HandleError(ctx, errors.ErrNotEnoughtQuota)
	}

	// We just check if the user exist for real
	if user.Username == "" {
		return errors.HandleError(ctx, errors.ErrDatabaseNotFound)
	}

	// Generate a random 32 bytes key for the file encryption
	fileKey, err := bcrypto.GenKey(chacha20poly1305.KeySize)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}

	// Get user's public key and transform it to string from byte data
	pbKey, err := ecies.NewPublicKeyFromBytes(user.PublicKey)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}

	// encrypt Xchacha20-Poly1305 key with user's public key
	encFileKey, err := ecies.Encrypt(pbKey, fileKey)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}

	// Create file from the model
	var file = models.File{
		Id:       uuid.New().String(),
		Owner:    token.Username,
		Name:     fileMultipart.Filename,
		Size:     fileMultipart.Size,
		Chunk:    0,
		Type:     fileMultipart.Header["Content-Type"][0],
		LastEdit: time.Now(),
	}

	// Create file access
	var access = models.Access{
		Id:            file.Id,
		State:         "PRIVATE",
		SharedBy:      token.Username,
		SharedTo:      token.Username,
		FileId:        file.Id,
		Favorite:      false,
		EncryptionKey: encFileKey,
	}

	// Get path storage path
	path := os.Getenv("STORAGE_DIR") + "/" + file.Id

	// Start piping the encrypted file to the server
	outfile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}
	defer outfile.Close()

	// Encrypting the file:
	// 	This function is cuting the file in 500 parts and encrypt each
	// 	of these parts with XChacha-Poly1305 and its 32 bytes randomly generated key
	//
	// Coming feature:
	// 	In the beta version, we will implement the federation protocole.
	// 	this part will probably modified to add the resilent backup system
	// 	which will send the parts to others instances of BLOC
	err = bcrypto.EncryptFile(*fileMultipart, fileKey, func(buf []byte, chunkSize int) {
		outfile.Write(buf) // write each parts to the output file
		// Put the size of the chunks in the database.
		// it's really important because during the download
		// XChacha-Poly1305 is verifying the integrity of the chunk
		// in order to decrypt it. And if the chunk is not of
		// the good size, the decryption will fail.
		file.Chunk = int64(chunkSize)
	})

	// If encryption fail
	if err != nil {
		return errors.HandleError(ctx, errors.ErrInternal)
	}

	// Add file metadatas to the database
	err = models.FileCreate(file)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrDatabaseCreate)
	}

	// Create access
	err = models.AccessCreate(access)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrDatabaseCreate)
	}

	err = models.UserUpdateQuota(user.Username, quota)
	if err != nil {
		return errors.HandleError(ctx, errors.ErrDatabaseUpdate)
	}

	// Respond with file metadatas
	return ctx.JSON(fiber.Map{
		"code":    "SUCCESS",
		"message": "File uploaded!",
		"data": fiber.Map{
			"file": fiber.Map{
				"accessId":    access.Id,
				"fileId":      file.Id,
				"accessState": access.State,
				"fileName":    file.Name,
				"fileType":    file.Type,
				"fileSize":    file.Size,
				"sharedBy":    access.SharedBy,
				"sharedTo":    access.SharedTo,
				"lastEdit":    file.LastEdit,
				"favorite":    access.Favorite,
			},
			"quota": fiber.Map{
				"max":   utils.GetQuota(),
				"total": quota,
			},
		},
	})
}
