package models

import (
	"github.com/herEmanuel/gocord/pkg/util/postgres"
)

type User struct {
	postgres.BaseModel
	Name     string
	Email    string   `json:",omitempty"`
	Password string   `json:",omitempty"`
	Avatar   string   `gorm:"default:-"`
	Servers  []Server `gorm:"many2many:user_servers" json:",omitempty"`
	Roles    []Role   `gorm:"many2many:roles_users" json:",omitempty"`
}
