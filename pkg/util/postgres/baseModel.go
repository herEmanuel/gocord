package postgres

import (
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	CreatedAt *time.Time `json:",omitempty"`
	UpdatedAt *time.Time `json:",omitempty"`
}

func (b *BaseModel) BeforeCreate(db *gorm.DB) error {
	userID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	b.ID = userID

	return nil
}
