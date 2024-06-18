package dto

import "mime/multipart"

type (
	JournalCreateRequest struct {
		Content string                `json:"content" form:"content" binding:"required"`
		Image   *multipart.FileHeader `json:"image" form:"image"`
		UserID  string                `json:"user_id" form:"user_id"`
	}

	JournalResponse struct {
		ID        string `json:"id"`
		Content   string `json:"content"`
		Image     string `json:"image,omitempty"`
		UserID    string `json:"user_id,omitempty"`
		Level     int    `json:"level,omitempty"`
		Point     int    `json:"point,omitempty"`
		CreatedAt string `json:"created_at"`
	}
)
