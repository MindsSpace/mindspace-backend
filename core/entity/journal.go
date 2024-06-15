package entity

import (
	"github.com/google/uuid"
	"github.com/zetsux/gin-gorm-clean-starter/common/base"
)

type Journal struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Content string    `json:"content" gorm:"not null"`
	Image   string    `json:"image"`
	UserID  string    `json:"user_id" gorm:"foreignKey:UserID;not null"`
	User    *User     `json:"user,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	base.Model
}
