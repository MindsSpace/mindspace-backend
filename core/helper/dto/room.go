package dto

import (
	"github.com/zetsux/gin-gorm-clean-starter/core/entity"
)

type (
	RoomCreateRequest struct {
		Name   string `json:"name" form:"name" binding:"required"`
		UserID string `json:"user_id" form:"user_id"`
	}

	RoomResponse struct {
		ID        string        `json:"id"`
		Name      string        `json:"name"`
		UserID    string        `json:"user_id,omitempty"`
		Chats     []entity.Chat `json:"chats,omitempty"`
		CreatedAt string        `json:"created_at"`
	}
)