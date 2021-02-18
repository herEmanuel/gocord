package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/herEmanuel/gocord/pkg/api/server/http"
	"github.com/herEmanuel/gocord/pkg/api/server/storage"
	"github.com/herEmanuel/gocord/pkg/util/auth"
	"gorm.io/gorm"
)

func Init(app *fiber.App, dbConn *gorm.DB) {

	storage.Db = dbConn
	v1 := app.Group("/v1")

	v1.Get("/getServer", auth.AuthMiddleware, http.GetServer)
	v1.Post("/createServer", auth.AuthMiddleware, http.CreateServer)
	v1.Post("/sendMessage/:channelID", auth.AuthMiddleware, http.SendMessage)
	v1.Get("/getChannelMessages/:channelID", auth.AuthMiddleware, http.GetChannelMessages)
}
