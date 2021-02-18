package models

import (
	"github.com/google/uuid"
	"github.com/herEmanuel/gocord/pkg/util/postgres"
)

type Channel struct {
	postgres.BaseModel
	Name       string
	Server     uuid.UUID
	Messages   []Message `gorm:"constraint:OnDelete:CASCADE;foreignKey:channel"`
	Permission string    //read-only, admin-only, public
}
