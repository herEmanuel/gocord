package models

import (
	"github.com/google/uuid"
	"github.com/herEmanuel/gocord/pkg/util/postgres"
)

type Role struct {
	postgres.BaseModel
	Name     string
	Color    string
	Priority uint8
	Server   uuid.UUID
	Users    []User `gorm:"many2many:roles_users"`
}
