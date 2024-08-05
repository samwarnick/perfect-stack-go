package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID        string    `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	Message   string    `gorm:"type:varchar(255);not null" json:"message,omitempty"`
}

func (message *Message) BeforeCreate(tx *gorm.DB) (err error) {
	message.ID = uuid.New().String()
	return nil
}

type CreateMessageSchema struct {
	Message string `json:"message" validate:"required"`
}
