package models

import (
	"github.com/herEmanuel/gocord/pkg/util/postgres"
)

type User struct {
	postgres.BaseModel
	Name     string
	Email    string
	Password string
	Avatar   string   `gorm:"default:-"`
	Servers  []Server `gorm:"many2many:user_servers"`
}
