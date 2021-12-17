package utils

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GenCookie(name string, value string, exp time.Duration, domain string) *fiber.Cookie {
	serverHTTPS, err := strconv.ParseBool(os.Getenv("SERVER_HTTPS"))
	if err != nil {
		log.Fatalln("Set SERVER_HTTPS as true or false please")
	}

	return &fiber.Cookie{
		Name:     name,
		Value:    value,
		SameSite: "strict",
		Secure:   serverHTTPS,
		HTTPOnly: true,
		Expires:  time.Now().Add(exp),
		Domain:   domain,
	}
}

func GetQuota() int {
	i, _ := strconv.Atoi(os.Getenv("STORAGE_QUOTA"))
	return i * 1024 * 1024
}
