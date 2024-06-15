package dto

import "mime/multipart"

type (
	UserAuthRequest struct {
		Username string `json:"username" form:"username" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	UserResponse struct {
		ID         string `json:"id"`
		Username   string `json:"username"`
		Level      int    `json:"level"`
		Point      int    `json:"point"`
		IsProfiled *bool  `json:"is_profiled,omitempty"`
		Avatar     string `json:"avatar,omitempty"`
	}

	UserUpdateRequest struct {
		ID       string `json:"id" form:"id"`
		Username string `json:"username" form:"username"`
		Level    int    `json:"level" form:"level"`
		Point    int    `json:"point" form:"point"`
		Password string `json:"password" form:"password"`
	}

	UserAddPointRequest struct {
		Point int `json:"point" form:"point" binding:"required"`
	}

	UserChangeAvatarRequest struct {
		Avatar *multipart.FileHeader `json:"avatar" form:"avatar" binding:"required"`
	}
)
