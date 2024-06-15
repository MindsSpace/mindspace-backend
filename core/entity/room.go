package entity

import (
	"github.com/google/uuid"
	"github.com/zetsux/gin-gorm-clean-starter/common/base"
)

type Room struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name   string    `json:"name" gorm:"not null"`
	Chats  []Chat    `json:"chats"`
	UserID string    `json:"user_id" gorm:"foreignKey:UserID;not null"`
	User   *User     `json:"user,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	base.Model
}
