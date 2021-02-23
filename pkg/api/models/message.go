package models

import (
	"github.com/google/uuid"
	"github.com/herEmanuel/gocord/pkg/util/postgres"
)

type Message struct {
	postgres.BaseModel
	Content string
	Type    string
	User    User      `json:",omitempty"`
	UserID  uuid.UUID `json:",omitempty"`
	Channel uuid.UUID `json:",omitempty"`
}
