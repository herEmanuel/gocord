package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/herEmanuel/gocord/pkg/api/user/http"
	"github.com/herEmanuel/gocord/pkg/api/user/storage"
	"gorm.io/gorm"
)

func Init(app *fiber.App, dbConn *gorm.DB) {

	storage.Db = dbConn
	v1 := app.Group("/v1")

	v1.Post("/register", http.Register)
	v1.Post("/login", http.Login)

}
