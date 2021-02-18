package models

import (
	"github.com/herEmanuel/gocord/pkg/util/postgres"
)

type Server struct {
	postgres.BaseModel
	Name       string
	Picture    string `gorm:"default:-"`
	InviteCode string
	Members    []User    `gorm:"many2many:user_servers"`
	Admins     []User    `gorm:"many2many:server_admins"`
	Channels   []Channel `gorm:"foreignKey:server;constraint:OnDelete:CASCADE"`
	Roles      []Role    `gorm:"foreignKey:server"`
}
