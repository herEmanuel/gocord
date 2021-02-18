package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/herEmanuel/gocord/pkg/api/server/gateway"
)

func CreateServer(ctx *fiber.Ctx) error {

	body := new(CreateServerJSON)

	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(500).SendString("Could not parse the request body")
	}

	userID := ctx.Locals("userID").(uuid.UUID)
	newServer, err := gateway.CreateServer(userID, body.Name)
	if err != nil {
		return ctx.Status(400).SendString("Could not create the new server, " + err.Error())
	}

	return ctx.JSON(map[string]interface{}{
		"ID":         newServer.ID,
		"Name":       newServer.Name,
		"InviteCode": newServer.InviteCode,
		"Picture":    newServer.Picture,
	})
}

func GetServer(ctx *fiber.Ctx) error {

	body := new(GetServerJSON)

	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(500).SendString("Could not parse the request body")
	}

	server, err := gateway.GetServer(body.ServerID)
	if err != nil {
		return ctx.Status(400).SendString("Could not find this server, " + err.Error())
	}

	return ctx.JSON(server)
}

func GetChannelMessages(ctx *fiber.Ctx) error {

	channelID, _ := uuid.Parse(ctx.Params("channelID"))
	userID := ctx.Locals("userID").(uuid.UUID)

	messages, err := gateway.GetChannelMessages(userID, channelID)
	if err != nil {
		return ctx.Status(400).SendString("Could not get the messages, " + err.Error())
	}

	return ctx.JSON(messages)
}

func SendMessage(ctx *fiber.Ctx) error {

	body := new(SendMessageJSON)

	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(500).SendString("Could not parse the request body")
	}

	channelID, _ := uuid.Parse(ctx.Params("channelID"))
	userID := ctx.Locals("userID").(uuid.UUID)

	err = gateway.SendMessage(userID, channelID, body.Content, body.MessageType)
	if err != nil {
		return ctx.Status(400).SendString("Could not send the message, " + err.Error())
	}

	return ctx.SendString("Your message was sent successfully")
}
