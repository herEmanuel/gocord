package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
		return ctx.Status(400).SendString("Could not register the new user, " + err.Error())
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
		return ctx.Status(400).SendString("Could not login, " + err.Error())
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "authToken",
		Value:    user["token"].(string),
		MaxAge:   60 * 60 * 24 * 7 * 1000,
		HTTPOnly: true,
	})
	return ctx.JSON(user["user"])
}

func ChangePassword(ctx *fiber.Ctx) error {

	body := new(ChangePasswordJSON)
	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(500).SendString("Could not parse the request body")
	}

	userID := ctx.Locals("userID").(uuid.UUID)

	err = gateway.ChangePassword(userID, body.OldPassword, body.NewPassword, body.ConfirmationPassword)
	if err != nil {
		return ctx.Status(400).SendString("Could not change your password, " + err.Error())
	}

	return ctx.SendString("Your password was changed successfully")
}

func AddImage(ctx *fiber.Ctx) error {

	userID := ctx.Locals("userID").(uuid.UUID)
	imagePath := ctx.Locals("imagePath").(string)

	err := gateway.AddImage(imagePath, userID)
	if err != nil {
		return ctx.Status(400).SendString("Could not save you new avatar, " + err.Error())
	}

	return ctx.SendString("Avatar updated successfully")
}

func EnterServer(ctx *fiber.Ctx) error {

	body := new(EnterServerJSON)
	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(500).SendString("Could not parse the request body")
	}

	userID := ctx.Locals("userID").(uuid.UUID)

	server, err := gateway.EnterServer(userID, body.InviteCode)
	if err != nil {
		return ctx.Status(400).SendString("Could not enter this server, " + err.Error())
	}

	return ctx.JSON(server)
}

func LeaveServer(ctx *fiber.Ctx) error {

	serverID, _ := uuid.Parse(ctx.Params("serverID"))
	userID := ctx.Locals("userID").(uuid.UUID)

	err := gateway.LeaveServer(userID, serverID)
	if err != nil {
		return ctx.Status(400).SendString("Could not leave the server, " + err.Error())
	}

	return ctx.SendString("You have left this server successfully")
}
