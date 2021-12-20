package users

import (
	"bloc/utils"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logout(ctx *fiber.Ctx) error {
	ctx.ClearCookie() // Remove all cookies

	// Remove all the cookies by setting their expiration time 30 seconds in the past
	ctx.Cookie(utils.GenCookie("token", "", time.Second*-30, os.Getenv("SERVER_DOMAIN")))
	if ctx.Hostname() == os.Getenv("TOR_ADDRESS") {
		ctx.Cookie(utils.GenCookie("token", "", time.Second*-30, os.Getenv("TOR_ADDRESS")))
	}

	return ctx.JSON(fiber.Map{
		"code": "SUCCESS",
	})
}
