package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/herEmanuel/gocord/pkg/api/server/chat"
	"github.com/herEmanuel/gocord/pkg/api/server/http"
	"github.com/herEmanuel/gocord/pkg/api/server/storage"
	"github.com/herEmanuel/gocord/pkg/util/auth"
	"github.com/herEmanuel/gocord/pkg/util/imageUpload"
	"gorm.io/gorm"
)

func Init(app *fiber.App, dbConn *gorm.DB) {

	storage.Db = dbConn
	chat.Init()
	app.Get("/ws", func(ctx *fiber.Ctx) error {
		if !websocket.IsWebSocketUpgrade(ctx) {
			return fiber.ErrUpgradeRequired
		}
		return ctx.Next()
	}, websocket.New(chat.WSConn))

	v1 := app.Group("/v1")

	v1.Post("/createServer", auth.AuthMiddleware, http.CreateServer)
	v1.Delete("/deleteServer/:serverID", auth.AuthMiddleware, auth.IsServerAdmin, http.DeleteServer)
	v1.Post("/addImage/:serverID", auth.AuthMiddleware, auth.IsServerAdmin, imageUpload.ImageUpload, http.AddImage)
	v1.Get("/getServer/:serverID", auth.AuthMiddleware, http.GetServer)
	v1.Post("/createChannel/:serverID", auth.AuthMiddleware, auth.IsServerAdmin, http.CreateChannel)
	v1.Delete("/deleteChannel/:serverID/:channelID", auth.AuthMiddleware, auth.IsServerAdmin, http.DeleteChannel)
	v1.Get("/getChannelMessages/:channelID", auth.AuthMiddleware, http.GetChannelMessages)
	v1.Post("/sendMessage/:channelID", auth.AuthMiddleware, http.SendMessage)
	v1.Delete("/deleteMessage/:messageID", auth.AuthMiddleware, http.DeleteMessage)
	v1.Post("/createRole/:serverID", auth.AuthMiddleware, auth.IsServerAdmin, http.CreateRole)
	v1.Delete("/deleteRole/:serverID/:roleID", auth.AuthMiddleware, auth.IsServerAdmin, http.DeleteRole)
	v1.Put("/addRoleToUser/:serverID", auth.AuthMiddleware, auth.IsServerAdmin, http.AddRoleToUser)
	v1.Delete("/removeUser/:serverID", auth.AuthMiddleware, auth.IsServerAdmin, http.RemoveUser)
	v1.Put("/promoteToAdmin/:serverID", auth.AuthMiddleware, auth.IsServerAdmin, http.PromoteToAdmin)
}
