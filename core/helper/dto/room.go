package dto

import (
	"github.com/zetsux/gin-gorm-clean-starter/core/entity"
)

type (
	RoomCreateRequest struct {
		UserID string `json:"user_id" form:"user_id"`
		Theme  string `json:"theme" form:"theme"`
	}

	RoomResponse struct {
		ID        string        `json:"id"`
		Name      string        `json:"name"`
		UserID    string        `json:"user_id,omitempty"`
		Chats     []entity.Chat `json:"chats,omitempty"`
		CreatedAt string        `json:"created_at"`
	}
)
