package dto

type (
	ChatCreateRequest struct {
		Content  string `json:"content" form:"content" binding:"required"`
		RoomID   string `json:"room_id" form:"room_id" binding:"required"`
		Language string `json:"language" form:"language" binding:"required"`
	}

	ChatResponse struct {
		ID        string `json:"id"`
		Content   string `json:"content"`
		IsUser    bool   `json:"is_user"`
		RoomID    string `json:"room_id"`
		CreatedAt string `json:"created_at"`
	}
)
