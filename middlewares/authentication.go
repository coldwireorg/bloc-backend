package middlewares

import (
	"bloc/utils/tokens"

	"github.com/gofiber/fiber/v2"
)

// Check if user have a correct token
func CheckUserToken(ctx *fiber.Ctx) error {
	token := ctx.Cookies("token")  // Get token cookie
	t, err := tokens.Verify(token) // Verify JWT token with public key
	if err != nil {
		// The token is not valid: we respond a 403 error code
		return ctx.Status(fiber.ErrForbidden.Code).JSON(fiber.Map{
			"code":    "ERROR_TOKEN",
			"message": "Invalid token",
		})
	} else {
		// We just chekck if the token is not empty (why not ?)
		if string(t.Token) != "" {
			return ctx.Next()
		} else {
			// if the cookie is empty, we "try" to clear it
			ctx.ClearCookie("token", "pvkey")
			return ctx.Redirect("/")
		}
	}
}
