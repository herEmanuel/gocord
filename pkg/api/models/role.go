package models

import (
	"github.com/herEmanuel/gocord/pkg/util/postgres"
)

type Role struct {
	postgres.BaseModel
	Name   string
	Color  string
	Server string
	Users  []User `gorm:"many2many:roles_users"`
}
