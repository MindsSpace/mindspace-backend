package entity

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/zetsux/gin-gorm-clean-starter/common/base"
)

type Profiling struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Mood       string         `json:"mood" gorm:"not null"`
	Problems   pq.StringArray `json:"problems" gorm:"type:varchar(100)[];not null"`
	Approaches pq.StringArray `json:"approaches" gorm:"type:varchar(100)[];not null"`
	UserID     string         `json:"user_id" gorm:"foreignKey:UserID;not null"`
	User       *User          `json:"user,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RoomID     string         `json:"room_id" gorm:"foreignKey:RoomID;not null"`
	Room       *Room          `json:"room,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	base.Model
}
