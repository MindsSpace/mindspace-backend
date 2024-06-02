package dto

type (
	UserAuthRequest struct {
		Username string `json:"username" form:"username" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	UserResponse struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Level    int    `json:"level"`
		Point    int    `json:"point"`
	}

	UserUpdateRequest struct {
		ID       string `json:"id" form:"id"`
		Username string `json:"username" form:"username"`
		Level    int    `json:"level" form:"level"`
		Point    int    `json:"point" form:"point"`
		Password string `json:"password" form:"password"`
	}
)
