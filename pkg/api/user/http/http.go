package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/herEmanuel/gocord/pkg/api/user/gateway"
)

func Register(ctx *fiber.Ctx) error {

	body := new(RegisterJSON)
	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(500).SendString("Could not parse the request body")
	}

	newUser, err := gateway.Register(body.Name, body.Email, body.Password, body.ConfirmPassword)
	if err != nil {
		return ctx.SendString("Could not register the new user, " + err.Error())
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "authToken",
		Value:    newUser["token"].(string),
		MaxAge:   60 * 60 * 24 * 7 * 1000,
		HTTPOnly: true,
	})
	return ctx.JSON(newUser["newUser"])
}

func Login(ctx *fiber.Ctx) error {

	body := new(LoginJSON)
	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(500).SendString("Could not parse the request body")
	}

	user, err := gateway.Login(body.Email, body.Password)
	if err != nil {
		return ctx.SendString("Could not login, " + err.Error())
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "authToken",
		Value:    user["token"].(string),
		MaxAge:   60 * 60 * 24 * 7 * 1000,
		HTTPOnly: true,
	})
	return ctx.JSON(user["user"])
}
