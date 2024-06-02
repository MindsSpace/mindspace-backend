package entity

import (
	"github.com/google/uuid"
	"github.com/zetsux/gin-gorm-clean-starter/common/base"
)

type Chat struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Content string    `json:"content" gorm:"not null"`
	IsUser  bool      `json:"is_user" gorm:"not null"`
	RoomID  string    `json:"room_id" gorm:"foreignKey:RoomID"`
	Room    *Room     `json:"room,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	base.Model
}
