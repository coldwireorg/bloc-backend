package routes

import (
	"bloc/controllers/favorites"
	"bloc/controllers/files"
	"bloc/controllers/shares"
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

	/* USERS RELATED ROUTES */
	user := api.Group("/user")
	user.Get("/quota", middlewares.CheckUserToken, users.QuotaGet)
	user.Post("/quota", middlewares.CheckUserToken, users.QuotaCheck)
	user.Post("/auth/login", users.Login)
	user.Post("/auth/register", users.Register)
	user.Post("/auth/logout", users.Logout)

	/* FILES RELATED ROUTES */
	file := api.Group("/file") // Route for servers
	file.Post("/", middlewares.CheckUserToken, files.Upload)
	file.Delete("/", middlewares.CheckUserToken, files.Delete)
	file.Get("/", middlewares.CheckUserToken, files.List)
	file.Post("/download", middlewares.CheckUserToken, files.Download)

	/* FAVORITES RELATED ROUTES */
	favorite := api.Group("/file") // Route for servers
	favorite.Post("/", middlewares.CheckUserToken, favorites.Update)

	/* SHARES RELATED ROUTES */
	share := api.Group("/share") // Route for servers
	share.Post("/", middlewares.CheckUserToken, shares.Create)
	share.Delete("/", middlewares.CheckUserToken, shares.Revoke)
	share.Get("/", middlewares.CheckUserToken, shares.List)
}
