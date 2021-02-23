package models

import (
	"github.com/google/uuid"
	"github.com/herEmanuel/gocord/pkg/util/postgres"
)

type Channel struct {
	postgres.BaseModel
	Name       string
	Permission string    //read-only, admin-only, public
	Server     uuid.UUID `json:",omitempty"`
	Messages   []Message `gorm:"constraint:OnDelete:CASCADE;foreignKey:channel" json:",omitempty"`
}
