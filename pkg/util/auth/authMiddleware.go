package auth

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/herEmanuel/gocord/pkg/util/jwt"
)

func AuthMiddleware(ctx *fiber.Ctx) error {

	authToken := ctx.Cookies("authToken")

	if len(authToken) != 0 {

		userID, err := jwt.VerifyToken(authToken, os.Getenv("JWT_SECRET_KEY"))
		if err != nil {
			return ctx.Status(401).SendString("Invalid token, " + err.Error())
		}

		ctx.Locals("userID", userID)
		return ctx.Next()
	}

	return ctx.Status(401).SendString("Invalid token")
}
