package routes

import (
	"bloc/controllers/files"
	"bloc/controllers/users"
	"bloc/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api") // Set route for API
	api.Use(logger.New())    // Add logger to the API routes

	// Check user token
	// This endpoint exist just for verifying auth cookies and the jwt token.
	// It simply makes the user's connection going through the auth middleware
	// and if the jwt is valid we answer a 200 code, if not the middleware will
	// answer a 403 (Forbidden) error code
	api.Get("/check", middlewares.CheckUserToken, func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{})
	})

	user := api.Group("/user") // Route for users
	file := api.Group("/file") // Route for servers

	file.Post("/quota", middlewares.CheckUserToken, files.CheckQuota)

	file.Post("/upload", middlewares.CheckUserToken, files.Upload)
	file.Delete("/delete", middlewares.CheckUserToken, files.Delete)

	file.Post("/favorite", middlewares.CheckUserToken, files.UpdateFavorite)

	file.Post("/share", middlewares.CheckUserToken, files.ShareFile)
	file.Delete("/share", middlewares.CheckUserToken, files.UnshareFile)

	file.Post("/download", middlewares.CheckUserToken, files.Download)
	file.Get("/list", middlewares.CheckUserToken, files.List)

	user.Post("/login", users.Login)
	user.Post("/register", users.Register)
	user.Post("/logout", users.Logout)
}
