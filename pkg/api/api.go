package api

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/herEmanuel/gocord/pkg/api/models"
	"github.com/herEmanuel/gocord/pkg/api/server"
	"github.com/herEmanuel/gocord/pkg/api/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Initialize() {

	db, err := gorm.Open(postgres.Open("host=localhost dbname="+os.Getenv("DB_NAME")+" port=5432 user=postgres password="+os.Getenv("DB_PASSWORD")), &gorm.Config{})
	if err != nil {
		panic("Error: could not connect to the database, " + err.Error())
	}

	db.AutoMigrate(
		&models.User{},
		&models.Server{},
		&models.Role{},
		&models.Message{},
		&models.Channel{},
	)

	app := fiber.New()

	user.Init(app, db)
	server.Init(app, db)

	app.Listen(":3000")
}