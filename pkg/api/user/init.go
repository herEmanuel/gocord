package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/herEmanuel/gocord/pkg/api/user/http"
	"github.com/herEmanuel/gocord/pkg/api/user/storage"
	"github.com/herEmanuel/gocord/pkg/util/auth"
	"github.com/herEmanuel/gocord/pkg/util/imageUpload"
	"gorm.io/gorm"
)

func Init(app *fiber.App, dbConn *gorm.DB) {

	storage.Db = dbConn
	v1 := app.Group("/v1")

	v1.Post("/register", http.Register)
	v1.Post("/login", http.Login)
	v1.Put("/changePassword", auth.AuthMiddleware, http.ChangePassword)
	v1.Put("/changeAvatar", auth.AuthMiddleware, imageUpload.ImageUpload, http.AddImage)
	v1.Get("/getUserInfo", auth.AuthMiddleware, http.GetUserInfo)
	v1.Put("/enterServer", auth.AuthMiddleware, http.EnterServer)
	v1.Put("/leaveServer/:serverID", auth.AuthMiddleware, http.LeaveServer)

}
