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

	return ctx.JSON(fiber.Map{
		"ID":         newServer.ID,
		"Name":       newServer.Name,
		"InviteCode": newServer.InviteCode,
		"Picture":    newServer.Picture,
	})
}

func DeleteServer(ctx *fiber.Ctx) error {

	serverID, _ := uuid.Parse(ctx.Params("serverID"))

	err := gateway.DeleteServer(serverID)
	if err != nil {
		return ctx.Status(400).SendString("Could not delete this server, " + err.Error())
	}

	return ctx.SendString("Server deleted successfully")
}

func AddImage(ctx *fiber.Ctx) error {

	serverID, _ := uuid.Parse(ctx.Params("serverID"))
	imagePath := ctx.Locals("imagePath").(string)

	err := gateway.AddImage(imagePath, serverID)
	if err != nil {
		return ctx.Status(400).SendString("Could not add the image to the server picture, " + err.Error())
	}

	return ctx.SendString("Image added successfully")
}

func GetServer(ctx *fiber.Ctx) error {

	serverID, _ := uuid.Parse(ctx.Params("serverID"))

	server, err := gateway.GetServer(serverID)
	if err != nil {
		return ctx.Status(400).SendString("Could not find this server, " + err.Error())
	}

	return ctx.JSON(server)
}

func CreateChannel(ctx *fiber.Ctx) error {

	body := new(CreateChannelJSON)

	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(500).SendString("Could not parse the request body")
	}

	serverID, _ := uuid.Parse(ctx.Params("serverID"))

	newChannel, err := gateway.CreateChannel(serverID, body.Name, body.Permission)
	if err != nil {
		return ctx.Status(400).SendString("Could not create the new channel, " + err.Error())
	}

	return ctx.JSON(newChannel)
}

func DeleteChannel(ctx *fiber.Ctx) error {

	channelID, _ := uuid.Parse(ctx.Params("channelID"))

	err := gateway.DeleteChannel(channelID)
	if err != nil {
		return ctx.Status(400).SendString("Could not delete this channel, " + err.Error())
	}

	return ctx.SendString("Channel deleted successfully")
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

	newMessage, err := gateway.SendMessage(userID, channelID, body.Content, body.MessageType)
	if err != nil {
		return ctx.Status(400).SendString("Could not send the message, " + err.Error())
	}

	return ctx.JSON(fiber.Map{
		"ID":        newMessage.ID,
		"CreatedAt": newMessage.CreatedAt,
		"Content":   newMessage.Content,
		"Type":      newMessage.Type,
		"ChannelID": newMessage.Channel,
	})
}

func DeleteMessage(ctx *fiber.Ctx) error {

	userID := ctx.Locals("userID").(uuid.UUID)
	messageID, _ := uuid.Parse(ctx.Params("messageID"))

	err := gateway.DeleteMessage(userID, messageID)
	if err != nil {
		return ctx.Status(400).SendString("Could not delete this message, " + err.Error())
	}

	return ctx.SendString("Message deleted successfully")
}

func CreateRole(ctx *fiber.Ctx) error {

	body := new(CreateRoleJSON)

	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(500).SendString("Could not parse the request body")
	}

	serverID, _ := uuid.Parse(ctx.Params("serverID"))

	role, err := gateway.CreateRole(serverID, body.Priority, body.Name, body.Color)
	if err != nil {
		return ctx.Status(400).SendString("Could not create the new role, " + err.Error())
	}

	return ctx.JSON(role)
}

func DeleteRole(ctx *fiber.Ctx) error {

	serverID, _ := uuid.Parse(ctx.Params("serverID"))
	roleID, _ := uuid.Parse(ctx.Params("roleID"))

	err := gateway.DeleteRole(roleID, serverID)
	if err != nil {
		return ctx.Status(400).SendString("Could not delete this role, " + err.Error())
	}

	return ctx.SendString("Role deleted successfully")
}

func AddRoleToUser(ctx *fiber.Ctx) error {

	body := new(AddRoleToUserJSON)

	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(500).SendString("Could not parse the request body")
	}

	roleID, _ := uuid.Parse(body.RoleID)
	userID, _ := uuid.Parse(body.UserID)
	serverID, _ := uuid.Parse(ctx.Params("serverID"))

	err = gateway.AddRoleToUser(roleID, userID, serverID)
	if err != nil {
		return ctx.Status(400).SendString("Could not add the role to this user, " + err.Error())
	}

	return ctx.SendString("Role added to this user successfully")
}

func RemoveUser(ctx *fiber.Ctx) error {

	body := new(RemoveUserJSON)

	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(500).SendString("Could not parse the request body")
	}

	serverID, _ := uuid.Parse(ctx.Params("serverID"))
	adminID := ctx.Locals("userID").(uuid.UUID)

	err = gateway.RemoveUser(body.UserID, adminID, serverID)
	if err != nil {
		return ctx.Status(400).SendString("Could not remove this user from this server, " + err.Error())
	}

	return ctx.SendString("User removed successfully")
}

func PromoteToAdmin(ctx *fiber.Ctx) error {

	body := new(PromoteToAdminJSON)

	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(500).SendString("Could not parse the request body")
	}

	serverID, _ := uuid.Parse(ctx.Params("serverID"))

	err = gateway.PromoteToAdmin(body.UserID, serverID)
	if err != nil {
		return ctx.Status(400).SendString("Could not promote this user to admin, " + err.Error())
	}

	return ctx.SendString("User promoted successfully")
}
