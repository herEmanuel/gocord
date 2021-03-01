package models

import (
	"github.com/herEmanuel/gocord/pkg/util/postgres"
)

type Server struct {
	postgres.BaseModel
	Name       string
	Picture    string `gorm:"default:-"`
	InviteCode string
	Members    []User    `gorm:"many2many:user_servers" json:",omitempty"`
	Admins     []User    `gorm:"many2many:server_admins" json:",omitempty"`
	Channels   []Channel `gorm:"foreignKey:server;constraint:OnDelete:CASCADE" json:",omitempty"`
	Roles      []Role    `gorm:"foreignKey:server;constraint:OnDelete:CASCADE" json:",omitempty"`
}
