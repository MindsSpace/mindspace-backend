package dto

type (
	ProfilingCreateRequest struct {
		Mood       string   `json:"mood" form:"mood" binding:"required"`
		Problems   []string `json:"problems" form:"problems" binding:"required"`
		Approaches []string `json:"approaches" form:"approaches" binding:"required"`
		UserID     string   `json:"user_id" form:"user_id"`
	}

	ProfilingResponse struct {
		ID         string   `json:"id"`
		Mood       string   `json:"mood"`
		Problems   []string `json:"problems"`
		Approaches []string `json:"approaches"`
		IsFilled   bool     `json:"is_filled"`
		UserID     string   `json:"user_id"`
		RoomID     string   `json:"room_id"`
		CreatedAt  string   `json:"created_at"`
	}
)
