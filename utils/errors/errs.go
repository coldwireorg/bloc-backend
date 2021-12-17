package errors

import "github.com/gofiber/fiber/v2"

var (
	/* INTERNAL */
	ErrRequest    = apiError{status: fiber.StatusBadRequest, code: "ERROR_REQUEST", message: "the server can't parse the datas sent by the client"}
	ErrInternal   = apiError{status: fiber.StatusInternalServerError, code: "ERROR_INTERNAL", message: "an internal server error occured"}
	ErrPermission = apiError{status: fiber.StatusForbidden, code: "ERROR_PERMISSION", message: "you do not have the required permision for this"}

	/* DATABASE */
	ErrDatabaseCreate       = apiError{status: fiber.StatusInternalServerError, code: "ERROR_DB", message: "failed to put data in the database"}
	ErrDatabaseUpdate       = apiError{status: fiber.StatusInternalServerError, code: "ERROR_DB", message: "failed update data in the database"}
	ErrDatabaseRemove       = apiError{status: fiber.StatusInternalServerError, code: "ERROR_DB", message: "failed remove put data in the database"}
	ErrDatabaseNotFound     = apiError{status: fiber.StatusNotFound, code: "ERROR_DB", message: "failed to find data in database"}
	ErrDatabaseAlreadyExist = apiError{status: fiber.StatusConflict, code: "ERROR_DB", message: "already exist"}

	/* AUTH */
	ErrAuth         = apiError{status: fiber.StatusForbidden, code: "ERROR_AUTH", message: "an issue occured during authentication"}
	ErrAuthExist    = apiError{status: fiber.StatusForbidden, code: "ERROR_AUTH_EXIST", message: "user already exist"}
	ErrAuthPassword = apiError{status: fiber.StatusForbidden, code: "ERROR_AUTH_PASSWORD", message: "wrong password"}

	/* QUOTA */
	ErrNotEnoughtQuota = apiError{status: fiber.StatusForbidden, code: "ERROR_QUOTA", message: "not enough space available"}
)
