package auth

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/herEmanuel/gocord/pkg/api/models"
	"github.com/herEmanuel/gocord/pkg/api/server/storage"
	"gorm.io/gorm"
)

func IsServerAdmin(ctx *fiber.Ctx) error {

	var server models.Server
	userID := ctx.Locals("userID").(uuid.UUID)
	serverID, err := uuid.Parse(ctx.Params("serverID"))
	if err != nil {
		return ctx.Status(500).SendString("Invalid server id, " + err.Error())
	}

	result := storage.Db.Preload("Admins").
		First(&server, "id = ?", serverID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ctx.Status(400).SendString("This server doesn't exist")
	}

	isAdmin := false
	for _, admin := range server.Admins {
		if admin.ID == userID {
			isAdmin = true
		}
	}

	if !isAdmin {
		return ctx.Status(401).SendString("You are not an admin in this server")
	}

	return ctx.Next()
}
